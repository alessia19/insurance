rm -r ~/.insuranceCLI
rm -r ~/.insuranceD
% update to wherever is the Go bin folder in your machine
PATH="$PATH:$HOME/go/bin"
insuranceD init mynode --chain-id insurance
insuranceCLI config keyring-backend test
insuranceCLI keys add oracle
insuranceCLI keys add client1
insuranceCLI keys add client2
insuranceD add-genesis-account $(insuranceCLI keys show oracle -a) 100000000foo,100000000stake
insuranceD add-genesis-account $(insuranceCLI keys show client1 -a) 100foo
insuranceD add-genesis-account $(insuranceCLI keys show client2 -a) 100foo
insuranceCLI config chain-id insurance
insuranceCLI config output json
insuranceCLI config indent true
insuranceCLI config trust-node true
insuranceD gentx --name oracle --keyring-backend test
insuranceD collect-gentxs