package test

import (
    "strings"
    "testing"

    "github.com/onflow/cadence"
    jsoncdc "github.com/onflow/cadence/encoding/json"
    emulator "github.com/onflow/flow-emulator"
    sdk "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
    sdktemplates "github.com/onflow/flow-go-sdk/templates"
    "github.com/onflow/flow-go-sdk/test"

    "github.com/stretchr/testify/assert"

    "github.com/onflow/flow-go-sdk"

    nft_contracts "github.com/onflow/flow-nft/lib/go/contracts"
)

const (
    commonNftRootPath                       = "../../.."
    commonNftCommonNftPath                  = commonNftRootPath + "/contracts/CommonNFT.cdc"
    commonNftSetupAccountPath               = commonNftRootPath + "/transactions/commonNFT/setup_account.cdc"
    commonNftMintCommonNftPath              = commonNftRootPath + "/transactions/commonNFT/mint_common_nft.cdc"
    commonNftMintMultipleCommonNftsPath     = commonNftRootPath + "/transactions/commonNFT/mint_multiple_common_nfts.cdc"
    commonNftTransferCommonNftPath          = commonNftRootPath + "/transactions/commonNFT/transfer_common_nft.cdc"
    commonNftTransferMultipleCommonNftsPath = commonNftRootPath + "/transactions/commonNFT/transfer_multiple_common_nfts.cdc"
    commonNftDeleteMultipleCommonNftsPath   = commonNftRootPath + "/transactions/commonNFT/delete_multiple_common_nfts.cdc"
    commonNftInspectCommonNftSupplyPath     = commonNftRootPath + "/scripts/commonNFT/read_common_nft_supply.cdc"
    commonNftInspectCollectionLenPath       = commonNftRootPath + "/scripts/commonNFT/read_collection_length.cdc"
    commonNftInspectCollectionIdsPath       = commonNftRootPath + "/scripts/commonNFT/read_collection_ids.cdc"
    commonNftInspectCollectionNftsPath      = commonNftRootPath + "/scripts/commonNFT/read_collection_nfts.cdc"
    commonNftInspectSingleNftPath           = commonNftRootPath + "/scripts/commonNFT/read_common_nft.cdc"
)

func CommonNftDeployContracts(b *emulator.Blockchain, t *testing.T) (flow.Address, flow.Address, crypto.Signer) {
    accountKeys := test.AccountKeyGenerator()

    // Should be able to deploy a contract as a new account with no keys.
    nftCode := loadNonFungibleToken()
    nftAddr, err := b.CreateAccount(
        nil,
        []sdktemplates.Contract{
            {
                Name:   "NonFungibleToken",
                Source: string(nftCode),
            },
        })
    if !assert.NoError(t, err) {
        t.Log(err.Error())
    }
    _, err = b.CommitBlock()
    assert.NoError(t, err)

    // Should be able to deploy a contract as a new account with one key.
    commonNftAccountKey, commonNftSigner := accountKeys.NewWithSigner()
    commonNftCode := loadCommonNft(nftAddr.String())
    commonNftAddr, err := b.CreateAccount(
        []*flow.AccountKey{commonNftAccountKey},
        []sdktemplates.Contract{
            {
                Name:   "CommonNFT",
                Source: string(commonNftCode),
            },
        })
    if !assert.NoError(t, err) {
        t.Log(err.Error())
    }
    _, err = b.CommitBlock()
    assert.NoError(t, err)

    // Simplify the workflow by having the contract address also be our initial test collection.
    CommonNftSetupAccount(t, b, commonNftAddr, commonNftSigner, nftAddr, commonNftAddr)

    return nftAddr, commonNftAddr, commonNftSigner
}

func CommonNftSetupAccount(t *testing.T, b *emulator.Blockchain, userAddress sdk.Address, userSigner crypto.Signer, nftAddr sdk.Address, commonNftAddr sdk.Address) {
    tx := flow.NewTransaction().
        SetScript(commonNftGenerateSetupAccountScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(100).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(userAddress)

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, userAddress},
        []crypto.Signer{b.ServiceKey().Signer(), userSigner},
        false,
    )
}

func CommonNftCreateAccount(t *testing.T, b *emulator.Blockchain, nftAddr sdk.Address, commonNftAddr sdk.Address) (sdk.Address, crypto.Signer) {
    userAddress, userSigner, _ := createAccount(t, b)
    CommonNftSetupAccount(t, b, userAddress, userSigner, nftAddr, commonNftAddr)
    return userAddress, userSigner
}

func CommonNftMintItem(b *emulator.Blockchain, t *testing.T, nftAddr, commonNftAddr flow.Address, commonNftSigner crypto.Signer) {
    tx := flow.NewTransaction().
        SetScript(commonNftGenerateMintCommonNftScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(100).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(commonNftAddr)
    tx.AddArgument(cadence.NewAddress(commonNftAddr))
    tx.AddArgument(cadence.NewUInt64(1))
    tx.AddArgument(cadence.NewString("DEVELOPER-METADATA"))
    tx.AddArgument(cadence.NewString("CONTENT-URL"))

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, commonNftAddr},
        []crypto.Signer{b.ServiceKey().Signer(), commonNftSigner},
        false,
    )
}

func CommonNftMintMultipleItems(b *emulator.Blockchain, t *testing.T, nftAddr, commonNftAddr flow.Address, commonNftSigner crypto.Signer, startEdition uint64, number uint64) {
    tx := flow.NewTransaction().
        SetScript(commonNftGenerateMintMultipleCommonNftsScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(200).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(commonNftAddr)
    tx.AddArgument(cadence.NewAddress(commonNftAddr))
    tx.AddArgument(cadence.NewUInt64(1))
    tx.AddArgument(cadence.NewUInt64(startEdition))
    tx.AddArgument(cadence.NewUInt64(number))
    tx.AddArgument(cadence.NewString("DEVELOPER-METADATA"))
    tx.AddArgument(cadence.NewString("CONTENT-URL"))

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, commonNftAddr},
        []crypto.Signer{b.ServiceKey().Signer(), commonNftSigner},
        false,
    )
}

func CommonNftTransferItem(b *emulator.Blockchain, t *testing.T, nftAddr, commonNftAddr flow.Address, commonNftSigner crypto.Signer, recipientAddr flow.Address, shouldFail bool) {
    tx := flow.NewTransaction().
        SetScript(commonNftGenerateTransferCommonNftScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(100).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(commonNftAddr)
    tx.AddArgument(cadence.NewAddress(recipientAddr))
    tx.AddArgument(cadence.NewUInt64(1))

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, commonNftAddr},
        []crypto.Signer{b.ServiceKey().Signer(), commonNftSigner},
        shouldFail,
    )
}

func CommonNftTransferMultipleItems(b *emulator.Blockchain, t *testing.T, nftAddr, commonNftAddr flow.Address, commonNftSigner crypto.Signer, recipientAddr flow.Address, shouldFail bool) {
    var nftIDArray []cadence.Value
    nftIDArray = append(nftIDArray, cadence.NewUInt64(1))
    nftIDArray = append(nftIDArray, cadence.NewUInt64(2))
    nftIDArray = append(nftIDArray, cadence.NewUInt64(3))

    tx := flow.NewTransaction().
        SetScript(commonNftGenerateTransferMultipleCommonNftsScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(100).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(commonNftAddr)
    tx.AddArgument(cadence.NewAddress(recipientAddr))
    tx.AddArgument(cadence.NewArray(nftIDArray))

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, commonNftAddr},
        []crypto.Signer{b.ServiceKey().Signer(), commonNftSigner},
        shouldFail,
    )
}

func CommonNftDeleteMultipleItems(b *emulator.Blockchain, t *testing.T, nftAddr, commonNftAddr flow.Address, commonNftSigner crypto.Signer, shouldFail bool) {
    var nftIDArray []cadence.Value
    nftIDArray = append(nftIDArray, cadence.NewUInt64(1))
    nftIDArray = append(nftIDArray, cadence.NewUInt64(2))
    nftIDArray = append(nftIDArray, cadence.NewUInt64(3))

    tx := flow.NewTransaction().
        SetScript(commonNftGenerateDeleteMultipleCommonNftsScript(nftAddr.String(), commonNftAddr.String())).
        SetGasLimit(100).
        SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
        SetPayer(b.ServiceKey().Address).
        AddAuthorizer(commonNftAddr)
    tx.AddArgument(cadence.NewArray(nftIDArray))

    signAndSubmit(
        t, b, tx,
        []flow.Address{b.ServiceKey().Address, commonNftAddr},
        []crypto.Signer{b.ServiceKey().Signer(), commonNftSigner},
        shouldFail,
    )
}

func TestCommonNftDeployContracts(t *testing.T) {
    b := newEmulator()
    CommonNftDeployContracts(b, t)
}

func TestCreateCommonNft(t *testing.T) {
    b := newEmulator()

    nftAddr, commonNftAddr, commonNftSigner := CommonNftDeployContracts(b, t)

    supply := executeScriptAndCheck(t, b, commonNftGenerateInspectCommonNftSupplyScript(nftAddr.String(), commonNftAddr.String()), nil)
    assert.Equal(t, cadence.NewUInt64(0), supply.(cadence.UInt64))

    len := executeScriptAndCheck(
        t,
        b,
        commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
        [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
    )
    assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))

    t.Run("Should be able to mint a CommonNFT", func(t *testing.T) {
        CommonNftMintItem(b, t, nftAddr, commonNftAddr, commonNftSigner)

        // Assert that the account's collection is correct
        len := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
        )
        assert.Equal(t, cadence.NewInt(1), len.(cadence.Int))

        // Retrieve the NFT
        nftRef := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectSingleNftScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr)), jsoncdc.MustEncode(cadence.NewUInt64(1))},
        )
        nftInfo := parseSingleNft(nftRef.(cadence.Optional))
        assert.Equal(t, uint64(1), nftInfo.ID)
        assert.Equal(t, uint64(1), nftInfo.DeveloperID)
        assert.Equal(t, uint64(0), nftInfo.Edition)
        assert.Equal(t, "DEVELOPER-METADATA", nftInfo.DeveloperMetadata)
        assert.Equal(t, "CONTENT-URL", nftInfo.ContentURL)
    })
}

func TestCreateMultipleCommonNfts(t *testing.T) {
    b := newEmulator()

    nftAddr, commonNftAddr, commonNftSigner := CommonNftDeployContracts(b, t)

    supply := executeScriptAndCheck(t, b, commonNftGenerateInspectCommonNftSupplyScript(nftAddr.String(), commonNftAddr.String()), nil)
    assert.Equal(t, cadence.NewUInt64(0), supply.(cadence.UInt64))

    len := executeScriptAndCheck(
        t,
        b,
        commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
        [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
    )
    assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))

    t.Run("Should be able to mint multiple CommonNFTs", func(t *testing.T) {
        CommonNftMintMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, 1, 5)

        // Assert that the account's collection is correct
        len := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
        )
        assert.Equal(t, cadence.NewInt(5), len.(cadence.Int))

        // Retrieve NFTs in the collection
        nftRefs := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionNftsScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
        )
        nftInfos := parseNfts(nftRefs.(cadence.Array))
        for i, info := range nftInfos {
            assert.Equal(t, uint64(i + 1), info.ID)
            assert.Equal(t, uint64(1), info.DeveloperID)
            assert.Equal(t, uint64(i + 1), info.Edition)
            assert.Equal(t, "DEVELOPER-METADATA", info.DeveloperMetadata)
            assert.Equal(t, "CONTENT-URL", info.ContentURL)
        }
    })
}

func TestTransferNft(t *testing.T) {
    b := newEmulator()

    nftAddr, commonNftAddr, commonNftSigner := CommonNftDeployContracts(b, t)

    userAddress, userSigner, _ := createAccount(t, b)

    // Create a new Collection
    t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
        CommonNftSetupAccount(t, b, userAddress, userSigner, nftAddr, commonNftAddr)

        len := executeScriptAndCheck(
            t,
            b, commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
        )
        assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))
    })

    // Transfer an NFT
    t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {
        CommonNftMintItem(b, t, nftAddr, commonNftAddr, commonNftSigner)

        CommonNftTransferItem(b, t, nftAddr, commonNftAddr, commonNftSigner, userAddress, false)

        // Assert that the account's collection is correct
        nftRef := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectSingleNftScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress)), jsoncdc.MustEncode(cadence.NewUInt64(1))},
        )
        nftInfo := parseSingleNft(nftRef.(cadence.Optional))
        assert.Equal(t, uint64(1), nftInfo.ID)
        assert.Equal(t, uint64(1), nftInfo.DeveloperID)
        assert.Equal(t, uint64(0), nftInfo.Edition)
        assert.Equal(t, "DEVELOPER-METADATA", nftInfo.DeveloperMetadata)
        assert.Equal(t, "CONTENT-URL", nftInfo.ContentURL)
    })

    // Transfer a non-existing NFT
    t.Run("Shouldn't be able to transer a non-existing NFT", func(t *testing.T) {
        CommonNftTransferItem(b, t, nftAddr, commonNftAddr, commonNftSigner, userAddress, true)
    })
}

func TestTransferMultipleNfts(t *testing.T) {
    b := newEmulator()

    nftAddr, commonNftAddr, commonNftSigner := CommonNftDeployContracts(b, t)

    userAddress, userSigner, _ := createAccount(t, b)

    // Create a new Collection
    t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
        CommonNftSetupAccount(t, b, userAddress, userSigner, nftAddr, commonNftAddr)

        len := executeScriptAndCheck(
            t,
            b, commonNftGenerateInspectCollectionLenScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
        )
        assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))
    })

    // Transfer multiple NFTs
    t.Run("Should be able to withdraw multiple NFTs and deposit to another accounts collection", func(t *testing.T) {
        CommonNftMintMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, 1, 5)

        CommonNftTransferMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, userAddress, false)

        // Assert that the account's collection is correct
        nftRefs := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionNftsScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
        )
        nftInfos := parseNfts(nftRefs.(cadence.Array))
        for i, info := range nftInfos {
            assert.Equal(t, uint64(i + 4), info.ID)
            assert.Equal(t, uint64(1), info.DeveloperID)
            assert.Equal(t, uint64(i + 4), info.Edition)
            assert.Equal(t, "DEVELOPER-METADATA", info.DeveloperMetadata)
            assert.Equal(t, "CONTENT-URL", info.ContentURL)
        }

        nftRefs2 := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionNftsScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
        )
        nftInfos2 := parseNfts(nftRefs2.(cadence.Array))
        for i, info := range nftInfos2 {
            assert.Equal(t, uint64(i + 1), info.ID)
            assert.Equal(t, uint64(1), info.DeveloperID)
            assert.Equal(t, uint64(i + 1), info.Edition)
            assert.Equal(t, "DEVELOPER-METADATA", info.DeveloperMetadata)
            assert.Equal(t, "CONTENT-URL", info.ContentURL)
        }
    })

    // Transfer multiple non-existing NFTs
    t.Run("Shouldn't be able to transer non-existing NFTs", func(t *testing.T) {
        CommonNftTransferMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, userAddress, true)
    })
}

func TestDeleteMultipleNfts(t *testing.T) {
    b := newEmulator()

    nftAddr, commonNftAddr, commonNftSigner := CommonNftDeployContracts(b, t)

    // Delete multiple NFTs
    t.Run("Should be able to delete multiple NFTs", func(t *testing.T) {
        CommonNftMintMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, 1, 5)

        CommonNftDeleteMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, false)

        // Assert that the account's collection is correct
        nftRefs := executeScriptAndCheck(
            t,
            b,
            commonNftGenerateInspectCollectionNftsScript(nftAddr.String(), commonNftAddr.String()),
            [][]byte{jsoncdc.MustEncode(cadence.NewAddress(commonNftAddr))},
        )
        nftInfos := parseNfts(nftRefs.(cadence.Array))
        for i, info := range nftInfos {
            assert.Equal(t, uint64(i + 4), info.ID)
            assert.Equal(t, uint64(1), info.DeveloperID)
            assert.Equal(t, uint64(i + 4), info.Edition)
            assert.Equal(t, "DEVELOPER-METADATA", info.DeveloperMetadata)
            assert.Equal(t, "CONTENT-URL", info.ContentURL)
        }
    })

    // Delete multiple non-existing NFTs
    t.Run("Shouldn't be able to delete non-existing NFTs", func(t *testing.T) {
        CommonNftDeleteMultipleItems(b, t, nftAddr, commonNftAddr, commonNftSigner, true)
    })
}

func replaceCommonNftAddressPlaceholders(code, nftAddress, commonNftAddress string) []byte {
    return []byte(replaceStrings(
        code,
        map[string]string{
            "\"../../contracts/NonFungibleToken.cdc\"": "0x" + nftAddress,
            "\"../../contracts/CommonNFT.cdc\"": "0x" + commonNftAddress,
        },
    ))
}

func loadNonFungibleToken() []byte {
    return nft_contracts.NonFungibleToken()
}

func loadCommonNft(nftAddr string) []byte {
    return []byte(strings.ReplaceAll(
        string(readFile(commonNftCommonNftPath)),
        "\"./NonFungibleToken.cdc\"",
        "0x"+nftAddr,
    ))
}

func commonNftGenerateSetupAccountScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftSetupAccountPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateMintCommonNftScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftMintCommonNftPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateMintMultipleCommonNftsScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftMintMultipleCommonNftsPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateTransferCommonNftScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftTransferCommonNftPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateTransferMultipleCommonNftsScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftTransferMultipleCommonNftsPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateDeleteMultipleCommonNftsScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftDeleteMultipleCommonNftsPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateInspectCommonNftSupplyScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftInspectCommonNftSupplyPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateInspectCollectionLenScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftInspectCollectionLenPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateInspectCollectionIdsScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftInspectCollectionIdsPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateInspectCollectionNftsScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftInspectCollectionNftsPath)),
        nftAddr,
        commonNftAddr,
    )
}

func commonNftGenerateInspectSingleNftScript(nftAddr, commonNftAddr string) []byte {
    return replaceCommonNftAddressPlaceholders(
        string(readFile(commonNftInspectSingleNftPath)),
        nftAddr,
        commonNftAddr,
    )
}
