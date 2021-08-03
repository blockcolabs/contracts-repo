import CommonNFT from "../../contracts/CommonNFT.cdc"

// This scripts returns the number of CommonNFTs currently in existence.

pub fun main(): UInt64 {
    return CommonNFT.totalSupply
}
