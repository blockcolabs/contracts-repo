import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import CommonNFT from "../../contracts/CommonNFT.cdc"

// This transaction transfers a Common NFT from one account to another.

transaction(recipient: Address, withdrawID: UInt64) {
    prepare(signer: AuthAccount) {
        // Get the recipients public account object.
        let recipient = getAccount(recipient)

        // Borrow a reference to the signer's NFT collection.
        let collectionRef = signer.borrow<&CommonNFT.Collection>(from: CommonNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // Borrow a public reference to the receivers collection.
        let depositRef = recipient.getCapability(CommonNFT.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!

        // Withdraw the NFT from the owner's collection.
        let nft <- collectionRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection.
        depositRef.deposit(token: <-nft)
    }
}
