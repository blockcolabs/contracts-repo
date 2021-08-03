import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This transction uses the NFTMinter resource to mint new NFTs.
//
// It must be run with the account that has the minter resource
// stored at path /storage/NFTMinter.

transaction(recipient: Address, developerID: UInt64, startEdition: UInt64, number: UInt64, developerMetadata: String, contentURL: String) {
    
    // Local variable for storing the minter reference.
    let minter: &CommonNFT.NFTMinter

    prepare(signer: AuthAccount) {
        // Borrow a reference to the NFTMinter resource in storage.
        self.minter = signer.borrow<&CommonNFT.NFTMinter>(from: CommonNFT.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
        // Get the public account object for the recipient.
        let recipient = getAccount(recipient)

        // Borrow the recipient's public NFT collection reference.
        let receiver = recipient
            .getCapability(CommonNFT.CollectionPublicPath)!
            .borrow<&{CommonNFT.CommonNFTCollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        // Mint the NFTs and deposit them to the recipient's collection.
        self.minter.mintMultipleNFTs(recipient: receiver, developerID: developerID, startEdition: startEdition, number: number, developerMetadata: developerMetadata, contentURL: contentURL)
    }
}
