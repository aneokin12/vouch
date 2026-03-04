console.log("Secret from Vouch:", process.env.TEST_SECRET);
if (process.env.TEST_SECRET === "test_secret_123") {
    console.log("Success! Process injection worked.");
} else {
    console.log("Error: Expected secret not found or incorrect.");
    process.exit(1);
}
