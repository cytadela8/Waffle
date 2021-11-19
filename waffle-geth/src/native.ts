import ffi from 'ffi-napi';
import {join} from 'path';
import type {TransactionRequest} from '@ethersproject/abstract-provider';
import {utils} from 'ethers';

export const library = ffi.Library(join(__dirname, '../go/build/wafflegeth.dylib'), {
  cgoCurrentMillis: ['int', []],
  getBlockNumber: ['string', []],
  sendTransaction: ['void', ['string']],
  getBalance: ['string', ['string']],
  call: ['string', ['string']],
  getTransactionCount: ['string', ['string']]
});

export function cgoCurrentMillis() {
  return library.cgoCurrentMillis();
}

export function getBlockNumber(): string {
  return library.getBlockNumber();
}

export function call(msg: TransactionRequest) {
  return library.call(JSON.stringify(msg))
}
export function getBalance(address: string): string {
  return library.getBalance(address);
}

export function sendTransaction(data: string): void {
  return library.sendTransaction(data);
}

export function getTransactionCount(address: string): string {
  return library.getTransactionCount(address);
}