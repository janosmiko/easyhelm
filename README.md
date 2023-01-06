# EasyHelm - Helm Chart generator

Generate helm charts from a simple configuration file.

Make it simple to create the necessary boilerplate for Kubernetes helm charts

## Usage

Use the `config.yaml` example to see the structure of the configuration.

```bash
$ easyhelm generate
```

## Customization

You can customize the chart by adding your own templates to the `input/chart` directory.

Eg: extract the content of the `assets/chart` directory to the `input/chart` and modify it. 
The generator will use these files instead of the built-in static files and generate the new helm charts based on
these custom files.