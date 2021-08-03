import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This transaction configures an account to hold Common NFTs.

transaction {
    prepare(signer: AuthAccount) {
        // If the account doesn't already have a collection.
        if signer.borrow<&CommonNFT.Collection>(from: CommonNFT.CollectionStoragePath) == nil {

            // Create a new empty collection.
            let collection <- CommonNFT.createEmptyCollection()
            
            // Save it to the account.
            signer.save(<-collection, to: CommonNFT.CollectionStoragePath)

            // Create a public capability for the collection.
            signer.link<&CommonNFT.Collection{NonFungibleToken.CollectionPublic, CommonNFT.CommonNFTCollectionPublic}>(CommonNFT.CollectionPublicPath, target: CommonNFT.CollectionStoragePath)
        }
    }
}
