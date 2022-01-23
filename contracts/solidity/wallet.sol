// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Allowance.sol";

contract SharedWallet is Allowance {
    event MoneySent(address indexed beneficiary, uint amount);
    event MoneyReceived(address indexed from, uint amount);

    function sendMoney(address payable _to, uint _amount) public {
        require(isOwner(), "Unauthorized to send money");
        require(address(this).balance >= _amount, "There are not enough founds");
        require(allowance[_to] >= _amount, "Assigned allowance is not enough to perform this transaction");
        reduceAllowance(_to, _amount);

        emit MoneySent(_to, _amount);
        _to.transfer(_amount);
    }

    function getContractBalance() public view onlyOwner returns(uint) {
        return address(this).balance;
    }

    function getBalance() public view returns(uint) {
        return address(msg.sender).balance;
    }

    function getBalanceOf(address _address) public view onlyOwner returns(uint) {
        return address(_address).balance;
    }

    receive() external payable {
        emit MoneyReceived(msg.sender, msg.value);
    }
}