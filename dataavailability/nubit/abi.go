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
				"name": "nubitHeight",
				"type": "bytes",
				"internalType": "bytes"
				},
				{
				"name": "commitment",
				"type": "bytes",
				"internalType": "bytes"
				}
			]
			}
		],
		"stateMutability": "pure"
	}
]`
