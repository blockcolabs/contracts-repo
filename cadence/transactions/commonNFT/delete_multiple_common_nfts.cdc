import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This transaction burns multiple Common NFTs from an account.

transaction(deleteIDs: [UInt64]) {
    prepare(signer: AuthAccount) {
        // Borrow a reference to the signer's NFT collection.
        let collectionRef = signer.borrow<&CommonNFT.Collection>(from: CommonNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // Burn the NFTs from the owner's collection.
        collectionRef.deleteMultiple(deleteIDs: deleteIDs)
    }
}
