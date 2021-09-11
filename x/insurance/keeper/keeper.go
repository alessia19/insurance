package keeper

import (
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"strconv"
	"time"
)

// Keeper of the insurance store
type Keeper struct {
	bankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper creates an insurance keeper
func NewKeeper(storeKey sdk.StoreKey, bankKeeper bank.Keeper, cdc *codec.Codec) Keeper {
	keeper := Keeper{
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		cdc:        cdc,
	}
	return keeper
}

const insuranceStorePrefix = ":insurance:"

// getInsuranceStoreKey
func getInsuranceStoreKey(ID string) []byte {
	return []byte(insuranceStorePrefix + ID)
}

func (keeper Keeper) getInsuraceByID(ctx sdk.Context, id string) (types.Insurance, error) {
	store := ctx.KVStore(keeper.storeKey)

	insuranceKey := getInsuranceStoreKey(id)
	if !store.Has(insuranceKey) {
		return types.Insurance{}, fmt.Errorf("insurance with ID %s not found", id)
	}

	var insurance types.Insurance
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(insuranceKey), &insurance)
	return insurance, nil
}

func (keeper Keeper) BuyInsurance(ctx sdk.Context, msg types.Insurance) error {

	store := ctx.KVStore(keeper.storeKey)

	//Controlla se l'ID di msg esiste già, in tal caso viene modificato l'ID del msg
	// in modo che siano ID differenti

	for store.Has(getInsuranceStoreKey(msg.ID)) {
		number, _ := strconv.Atoi(msg.ID)
		number = number + 1
		ID := strconv.Itoa(number)

		// il ciclo serve per la lunghezza della stringa di ID, lunga 4
		for i := 0; i < 5; i++ {
			if len(ID) != 4 {
				ID = "0" + ID
			}
		}
		msg.ID = ID
	}

	if !store.Has(getInsuranceStoreKey(msg.ID)) {
		store.Set(getInsuranceStoreKey(msg.ID), keeper.cdc.MustMarshalBinaryBare(&msg))
		return nil
	}

	_, err := keeper.getInsuraceByID(ctx, msg.ID)

	if err != nil {
		return err
	}

	return nil
}

func (keeper Keeper) insuranceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(insuranceStorePrefix))
}

type insuranceDiscrimination func(insurance types.Insurance) bool

func (keeper Keeper) getInsurance(ctx sdk.Context, logic insuranceDiscrimination) []types.Insurance {
	ri := keeper.insuranceIterator(ctx)
	defer ri.Close()

	receivedResult := []types.Insurance{}
	for ; ri.Valid(); ri.Next() {
		var insurance types.Insurance
		keeper.cdc.MustUnmarshalBinaryBare(ri.Value(), &insurance)

		// solo assicurazioni che soddisfano il filtro
		if logic(insurance) {
			receivedResult = append(receivedResult, insurance)
		}
	}

	return receivedResult
}

func (keeper Keeper) GetAllInsurance(ctx sdk.Context) []types.Insurance {
	return keeper.getInsurance(ctx, func(_ types.Insurance) bool {
		return true
	})
}

// today restituisce la data e l'ora di oggi
func today() time.Time {
	return time.Now()
}

// isToday controlla che day, month, year passato coincida con il giorno, mese e anno di oggi
func isToday(day, month, year int) bool {
	return day == today().Day() && time.Month(month) == today().Month() && year == today().Year()
}

// rimuovi serve per eliminare le assicurazioni che non coincidano con la data odierna
func (keeper Keeper) rimuovi(ctx sdk.Context) (result []types.Insurance) {
	insurance := keeper.GetAllInsurance(ctx)

	for i := 0; i < len(insurance); i++ {
		data := time.Date(insurance[i].Year, time.Month(insurance[i].Month), insurance[i].Day, 0, 0, 0, 0, time.Local)

		// se coincide con la data di oggi viene aggiunto a result che andrà ad aggiungere tutte le insurance di oggi
		if isToday(data.Day(), types.MonthInt(data.Month()), data.Year()) {
			result = append(result, insurance[i])
		}
	}

	return result
}

// RepayInsurance utilizzato per rimborsare chi ha un'assicurazione contro il maltempo per la data di oggi
func (keeper Keeper) RepayInsurance(ctx sdk.Context, msg types.MsgRepayInsurance) error {
	insurance := keeper.rimuovi(ctx)

	if len(insurance) == 0 {
		return fmt.Errorf("no insurance as of today's date")
	}

	for i := 0; i < len(insurance); i++ {
		season := types.Now()                         // restituisce la stagione corrispondente
		importo := insurance[i].Indemnization(season) // importo da rimborsare

		if msg.Oracle.Equals(insurance[i].Intestatario) {
			return fmt.Errorf("the insurance cannot be made out to the oracle")
		}

		if err := keeper.bankKeeper.SendCoins(ctx, msg.Oracle, insurance[i].Intestatario, sdk.NewCoins(sdk.NewCoin(insurance[i].Money.Denom, sdk.NewInt(importo)))); err != nil {
			return err
		}

		insurance[i].Rimborso = insurance[i].Rimborso.Add(sdk.NewCoin(insurance[i].Money.Denom, sdk.NewInt(importo)))
		keeper.updateInsurance(ctx, insurance[i])
	}

	return nil

}

func (keeper Keeper) updateInsurance(ctx sdk.Context, insurance types.Insurance) error {

	store := ctx.KVStore(keeper.storeKey)

	if !store.Has(getInsuranceStoreKey(insurance.ID)) {
		return fmt.Errorf("insurance with %s not found", insurance.ID)
	}

	store.Set(getInsuranceStoreKey(insurance.ID), keeper.cdc.MustMarshalBinaryBare(&insurance))
	return nil
}
