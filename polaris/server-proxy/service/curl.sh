curl --location 'http://9.134.117.127:8099/naming/v1/services/metadata' \
--header 'Platform-Token: a63acf6a46fd44f1ad892f80a2332c13' \
--header 'Platform-Id: polaris-sdk-test' \
--header 'Content-Type: application/json' \
--header 'Cookie: x-client-ssid=74602415:019fac7c4aac:06053d; x_host_key_access=97220098170d203108424eb8f2eb97e750352040_s' \
--data '[
    {
    "name": "lzb_test3",
    "namespace": "Test",
    "token": "118805345d5942689a4953063dd94df1",
    "metadata": {
        "key": "value",
        "key1": "value2"
        }
    }
]'