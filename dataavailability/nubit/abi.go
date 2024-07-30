package nubit

const blobDataABI = `[
	{
		"type": "function",
		"name": "BlobData",
		"inputs": [
			{
			"name": "blobData",
			"type": "tuple",
			"internalType": "struct NubitDAVerifier.BlobData",
			"components": [
				{
				"name": "blobID",
				"type": "bytes",
				"internalType": "bytes"
				},
				{
				"name": "signature",
				"type": "bytes",
				"internalType": "bytes"
				}
			]
			}
		],
		"stateMutability": "pure"
	}
]`
