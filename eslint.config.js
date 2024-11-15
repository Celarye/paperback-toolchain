import eslint from "@eslint/js";
import tseslint from "typescript-eslint";

export default tseslint.config(
  {
    extends: [eslint.configs.recommended],
    files: ["**/*{.js,.ts}"],
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      allowBitwiseExpressions: true,
    },
  },
  {
    extends: [...tseslint.configs.strictTypeChecked],
    files: ["**/*.ts"],
    languageOptions: {
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    ignores: ["**/lib", "**/src/devtools/generated", "**/dist"],
  },
);
