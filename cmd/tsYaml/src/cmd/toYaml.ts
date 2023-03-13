
// convert TS file in to YAML.
// toYaml <TS file>
import * as YAML from "yaml"

async function main() {
    const args = process.argv.slice(2)
    if (args.length != 1) {
        console.error("Expected 1 argument. <TS file>")
        process.exit(1);
    }

    const ts = await import(args[0]);
    process.stdout.write(YAML.stringify(ts["default"], {version: "1.1"}))
}

main()