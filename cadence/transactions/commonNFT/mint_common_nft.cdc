import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This transction uses the NFTMinter resource to mint a new NFT.
//
// It must be run with the account that has the minter resource
// stored at path /storage/NFTMinter.

transaction(recipient: Address, developerID: UInt64, developerMetadata: String, contentURL: String) {
    
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
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        // Mint the NFT and deposit it to the recipient's collection.
        self.minter.mintNFT(recipient: receiver, developerID: developerID, developerMetadata: developerMetadata, contentURL: contentURL)
    }
}