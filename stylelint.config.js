export default {
  extends: ['stylelint-config-standard'],
  rules: {
    // Allow single-line blocks for compact CSS
    'declaration-block-single-line-max-declarations': null,
    // Don't require specific casing for custom properties
    'custom-property-pattern': null,
    // Allow vendor prefixes
    'property-no-vendor-prefix': null,
    'value-no-vendor-prefix': null,
    // Relax selector rules for existing code
    'selector-class-pattern': null,
    'selector-id-pattern': null,
    // Allow empty blocks (used for placeholder styling)
    'block-no-empty': null,
    // Disable empty line requirements (let Prettier handle formatting)
    'rule-empty-line-before': null,
    'at-rule-empty-line-before': null,
    'declaration-empty-line-before': null,
    // Allow camelCase keyframe names (used in existing code)
    'keyframes-name-pattern': null,
    // Allow legacy color notation (rgba, hsla)
    'color-function-notation': null,
    'alpha-value-notation': null,
    // Allow long hex colors
    'color-hex-length': null,
    // Allow currentColor casing
    'value-keyword-case': null,
    // CSS is organized by component, not specificity - selector order is intentional
    'no-descending-specificity': null,
    // Prettier handles comment formatting
    'comment-empty-line-before': null,
  },
}
