import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This script returns the metadata for an NFT in an account's collection.

pub fun main(address: Address, itemID: UInt64): &CommonNFT.NFT {

    // Get the public account object for the token owner.
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(CommonNFT.CollectionPublicPath)!
        .borrow<&{CommonNFT.CommonNFTCollectionPublic}>()
        ?? panic("Could not borrow CommonNFTCollectionPublic")

    // Borrow a reference to a specific NFT in the collection.
    let commonNFT = collectionBorrow.borrowCommonNFT(id: itemID)
        ?? panic("No such itemID in that collection")

    return commonNFT
}
