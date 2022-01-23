// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../vendor/openzeppelin/contracts/access/Ownable.sol";

contract Allowance is Ownable {

    event AllowanceChanged(address indexed beneficiary, address indexed sender, uint prevAmount, uint newAmount);

    mapping(address => uint) public allowance;

    function isOwner() internal view returns(bool) {
        return owner() == msg.sender;
    }

    function setAllowance(address _beneficiary, uint _amount) public onlyOwner {
        emit AllowanceChanged(_beneficiary, msg.sender, allowance[_beneficiary], _amount);
        allowance[_beneficiary] = _amount;
    }

    function increaseAllowance(address _beneficiary, uint _amount) public onlyOwner {
        emit AllowanceChanged(_beneficiary, msg.sender, allowance[_beneficiary], allowance[_beneficiary] + _amount);
        allowance[_beneficiary] += _amount;
    }

    function reduceAllowance(address _beneficiary, uint _amount) public onlyOwner {
        emit AllowanceChanged(_beneficiary, msg.sender, allowance[_beneficiary], allowance[_beneficiary] - _amount);
        allowance[_beneficiary] -= _amount;
    }
}