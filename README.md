# Shared Wallet with allowance
Solidity smart contract that handle allowance control for a distributed budget.

## Contract
The contract files are at `./contracts` folder. Third-party contracts like Open Zeppelin are stored at `./contracts/vendor` folder,
while compiled contracts are stored at `./contracts/bin` folder

### Compile contracts
In order to compile the contracts the `ethereum/solc` docker image is used as following:

```shell
docker run --rm -v /Users/steven/code/eth/sharedWallet/contracts:/contracts ethereum/solc:stable \
/contracts/solidity/wallet.sol --bin --abi --optimize --overwrite -o /contracts/bin
```

To generate a Go interface and bindings based on the contracts, the `ethereum/client-go:alltools-latest` docker image is used as following:

```shell
docker run --rm -v /Users/steven/code/eth/sharedWallet/contracts:/contracts ethereum/client-go:alltools-latest \
abigen --abi /contracts/bin/SharedWallet.abi --pkg contracts --type Contract --out \
/contracts/interfaces/shared_wallet.go  --bin /contracts/bin/SharedWallet.bin
```

## Geth

Run Geth docker image as:
```shell
docker run --rm --name ethereum-node -v /Users/steven/code/eth/:/root -p 8545:8545 -p 30303:30303 ethereum/client-go --nodiscover
```

and connect to the instance as following in order to perform CLI operations:
```shell
docker exec -it <> /bin/sh

$ geth attach /root/.ethereum/geth.ipc
```

If you like to have a GUI app instead of CLI you can use Ganache

## DAPP - Decentralized application
The Shared Wallet With Allowance DAPP is implemented as a CLI application which has a set of commands to deploy, monitor and run the contract operations

### Configuration
There are two main configurations:
* `blockchain` which contains the HTTP and WebSocket addresses of the blockchain (could be Ganache for development) and the private key of the account that will deploy the contract
* `contract` which contains the deployed address contract and gas information to perform the transactions in the blockchain

### Deploy
In order to deploy the contact the private key of the owner account should be set either in the configuration file, env variable or as `-k` flag. Then run `./wallet deploy`. 
The command will return the contract address that should be used to monitor and run the contract transactions. It could be set in the config file, environment or flag

### Monitor
In order to monitor the events generated in by the deployed contract just run the command `./wallet monitor`. If the contract address in incorrect you will get this message `invalid contract address`,
otherwise you will get the message `start monitoring at 0xdB1ad94CBFA75951ec5DCeAC4C7e829A58b0BfF1` and then JSON data for every event that is generated while executing the transactions, for example:
```json
{
  "event_type": "MoneyReceived",
  "sender": "0x3F9CD35D0159d961039780A48DB6b7c595D40Fa4",
  "block_number": 7,
  "amount": 50000000,
  "timestamp": "2022-01-17T18:36:06.635701-04:00"
}
```
### Commands
There are four group of commands that could be performed:
* Allowance: handle the allowance for a given beneficiary
* Balance: get the balance of the contract and beneficiaries
* Ownership: get the current owner or change the ownership
* Transfers receive ether in the contract and send ether to beneficiaries

#### Allowance
The base command is `./wallet run allowance --action=set --amount=1000 -t 0x3F` where `--action` flag could be `set`, `get`, `increase` or `reduce`.
The `--amount` flag is amount in `wei` that is set in the allowance for the given `-t` target address (beneficiary).

When the action `get` is used, the `--amount` flag is not required since we're getting the allowance set to the beneficiary.

`set`, `increase` and `reduce` actions will generate an `AllowanceChanged` event and the `monitor` process will display a JSON message like this:
```json
{
  "event_type": "AllowanceChanged",
  "sender": "0xC109730b9D3A84B49b4D455F6C75f94DA1417DcF",
  "beneficiary": "0x130323C2A1a2A5Ac385D85545E03DF2b4C57bbc5",
  "prev_amount": 1000,
  "new_amount": 1100,
  "timestamp": "2022-01-18T20:49:41.309777-04:00"
}
```

#### Balance
The base command is `./wallet run balanace` and it has 2 flags, `--of` which could be `contract` or `address`, and the `-t` target flag.

The `./wallet run balanace --of=contract` command returns the contract balance in `wei`

The `./wallet run balanace --of=address -t 0x3F` command returns the address balance in `wei`

#### Ownership
The base command is `./wallet run ownership` and it also has the `--action` and `-t` flags.

The `./wallet run ownership --action=get` command returns the current owner address

The `./wallet run ownership --action=transfer -t 0x3F` command transfer the ownership to the given address

#### Transfer
In order to receive ether in the contract (from the owner) run `./wallet run transfer --action=receive --amount=10`, the amount is in ether

To send ether to a beneficiary use `./wallet run transfer --action=send --amount=5 -t 0x5A` but make sure the beneficiary has allowance set