#!/bin/bash

cd ..

##### Set up project and accounts. #####

flow project deploy --update

flow accounts create --key 65dcc6b9f5bbfc8247ac11081eef721bec399fcbfb6ccba2575579b2fe5c78396f4bdc022378ab8e06a590c97e1d4981691127922535a6f30e87c0d070ae680f --signer emulator-account

flow transactions send cadence/transactions/flow/transfer_tokens.cdc --signer emulator-account --arg UFix64:100.0 --arg Address:01cf0e2f2f715450

flow transactions send cadence/transactions/commonNFT/setup_account.cdc --signer emulator-account

flow transactions send cadence/transactions/commonNFT/setup_account.cdc --signer emulator-account-2

##### Mint and trasnfer NFTs. #####

flow transactions send cadence/transactions/commonNFT/mint_common_nft.cdc --signer emulator-account --arg Address:01cf0e2f2f715450 --arg UInt64:1 --arg String:"Single NFT" --arg String:"single-nft"

flow transactions send cadence/transactions/commonNFT/mint_multiple_common_nfts.cdc --signer emulator-account --arg Address:01cf0e2f2f715450 --arg UInt64:1 --arg UInt64:1 --arg UInt64:5 --arg String:"Multiple NFTs" --arg String:"multiple-nfts"

flow transactions send cadence/transactions/commonNFT/transfer_common_nft.cdc --signer emulator-account-2 --arg Address:f8d6e0586b0a20c7 --arg UInt64:5

# Is there a way to specify an array as an argument?
# flow transactions send cadence/transactions/commonNFT/transfer_multiple_common_nfts.cdc --signer emulator-account-2 --arg Address:f8d6e0586b0a20c7 --arg [UInt64; 2]:[4, 6]

##### Read back to verify states. #####

flow scripts execute cadence/scripts/commonNFT/read_common_nft_supply.cdc

flow scripts execute cadence/scripts/commonNFT/read_collection_ids.cdc --arg Address:01cf0e2f2f715450

flow scripts execute cadence/scripts/commonNFT/read_collection_length.cdc --arg Address:01cf0e2f2f715450

flow scripts execute cadence/scripts/commonNFT/read_collection_ids.cdc --arg Address:f8d6e0586b0a20c7

flow scripts execute cadence/scripts/commonNFT/read_collection_length.cdc --arg Address:f8d6e0586b0a20c7

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:01cf0e2f2f715450 --arg UInt64:1

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:01cf0e2f2f715450 --arg UInt64:2

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:01cf0e2f2f715450 --arg UInt64:3

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:01cf0e2f2f715450 --arg UInt64:4

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:f8d6e0586b0a20c7 --arg UInt64:5

flow scripts execute cadence/scripts/commonNFT/read_common_nft.cdc --arg Address:01cf0e2f2f715450 --arg UInt64:6

flow scripts execute cadence/scripts/commonNFT/read_collection_nfts.cdc --arg Address:f8d6e0586b0a20c7

flow scripts execute cadence/scripts/commonNFT/read_collection_nfts.cdc --arg Address:01cf0e2f2f715450

##### Burn NFTs. #####

# Is there a way to specify an array as an argument?
#flow transactions send cadence/transactions/commonNFT/delete_multiple_common_nfts.cdc --signer emulator-account --arg [UInt64; 1]:[5]

#flow transactions send cadence/transactions/commonNFT/delete_multiple_common_nfts.cdc --signer emulator-account-2 --arg [UInt64; 5]:[1, 2, 3, 4, 6]
