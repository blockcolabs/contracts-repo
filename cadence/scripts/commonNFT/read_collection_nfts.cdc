import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This script returns the metadata for all the NFTs in an account's collection.

pub fun main(address: Address): [&NonFungibleToken.NFT] {

    // Get the public account object for the token owner.
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(CommonNFT.CollectionPublicPath)!
        .borrow<&{CommonNFT.CommonNFTCollectionPublic}>()
        ?? panic("Could not borrow CommonNFTCollectionPublic")

    // Borrow an array of references to all the NFTs in the collection.
    let commonNFTs = collectionBorrow.borrowAllCommonNFTs()

    return commonNFTs
}
