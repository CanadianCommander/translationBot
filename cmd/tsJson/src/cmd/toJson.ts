
// convert TS file to json.
// toJson <TS file>
async function main() {
    const args = process.argv.slice(2)
    if (args.length != 1) {
        console.error("Expected 1 argument. <TS file>")
        process.exit(1);
    }

    const ts = await import(`../../${args[0]}`);
    process.stdout.write(JSON.stringify(ts["default"]))
}

main()