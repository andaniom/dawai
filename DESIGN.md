# Design System: DAWAI — Instrument-Grade Precision
**Project ID:** projects/13385464508616838578 (Stitch)

## 1. Visual Theme & Atmosphere

The design system is built on a narrative of **Instrument-Grade Precision**, blending the warmth of a violin's spruce body with the exacting clarity of modern music notation. It is a professional tool designed for music educators operating in high-focus classroom environments, where the interface must be as reliable and legible as a well-printed score.

The visual style is **Corporate / Modern** with a **Minimalist** soul, emphasizing:
- **Warmth and Focus:** A move away from clinical, cold grays in favor of parchment-like surfaces and deep, ink-like typography.
- **Intentional Asymmetry:** Avoiding the generic "three-column card" layout for more dynamic, grounded compositions that feel deliberate rather than templated.
- **Professional Restraint:** Motion is functional, confirming actions without distraction. Decoration is eschewed in favor of structure and clarity.
- **Asymmetric Balance:** A high variance in layout that prioritizes information hierarchy over perfect symmetry.

## 2. Color Palette & Roles

The palette is strictly limited to ensure a focused professional atmosphere. **Parchment** serves as the canvas, evoking the feel of quality sheet music paper. **Rosewood** is the singular chromatic accent, used exclusively for interactive triggers, status highlights, and level badges.

| Name | Hex / Value | Role |
|---|---|---|
| Parchment | `#FAF8F5` | Primary canvas — main background |
| Ivory Surface | `#FFFFFF` | Card/container background — subtle lift off canvas |
| Ink Deep | `#1C1917` | Mandatory primary text color |
| Graphite | `#57534E` | Secondary metadata, placeholders |
| Rosewood | `#B5603C` | Singular accent — CTAs, focus states, level badges. Never gradients or neon variants |
| Rosewood Muted | `rgba(181, 96, 60, 0.10)` | Hover backgrounds, row selection |
| Whisper Border | `rgba(214, 211, 209, 0.5)` | All UI outlines |
| Caution Amber | `#D97706` | Offline banner, warning utility messaging |
| Destructive Clay | `#DC2626` | Error states (text + border) |
| Success Sage | `#16A34A` | Sync-success confirmation, auto-dismiss |
| Bronze Text / BG | `#92400E` / `rgba(146,64,14,0.15)` | Bronze level badge |
| Silver Text / BG | `#6B7280` / `rgba(107,114,128,0.15)` | Silver level badge |
| Gold Text / BG | `#B45309` / `rgba(217,119,6,0.15)` | Gold level badge |

Status colors are used sparingly for utility messaging (Offline, Sync, Error) — never as decoration.

**Dark mode** (from CLAUDE.md §13.7): Background `#1C1917`, Surface `#292524`, Text primary `#FAF8F5`, Text secondary `#A8A29E`, Border `rgba(87,83,78,0.5)`. Rosewood accent is unchanged in dark mode.

## 3. Typography Rules

Typography establishes hierarchy through weight contrast and family shifts rather than excessive size variation.

- **UI & Content:** `Satoshi` for all interface labels and descriptive text. Headings use tight tracking (`-0.02em`) for a crisp, editorial feel.
- **Data & Numbers:** `Geist Mono` is mandatory for all numerical data, timestamps, and score totals — prevents layout jitter during real-time assessment and reinforces the "precision tool" aesthetic.
- **Constraints:** Body text line length never exceeds `65ch`.
- **Banned:** Generic system fonts (Inter, Roboto) are strictly prohibited.

| Token | Family | Size | Weight | Line height | Notes |
|---|---|---|---|---|---|
| display-lg | Satoshi | 1.875rem | 600 | 2.25rem | `-0.02em` tracking |
| display-lg-mobile | Satoshi | 1.5rem | 600 | 1.75rem | `-0.02em` tracking |
| section-heading | Satoshi | 1.25rem | 500 | 1.75rem | |
| body-base | Satoshi | 1rem | 400 | 1.6rem | |
| metadata-sm | Satoshi | 0.875rem | 400 | 1.25rem | |
| score-display | Geist Mono | 3rem | 700 | 1 | |
| score-numeral | Geist Mono | 2rem | 700 | 1 | |
| grade-badge | Geist Mono | 0.75rem | 600 | 1 | |
| level-badge-text | Geist Mono | 0.7rem | 500 | — | |

## 4. Component Stylings

**Buttons**
- Primary: Rosewood fill, Ivory text. Active state: `-1px` Y-axis translate + subtle darken. No glow.
- Secondary: Whisper-border outline, Ink Deep text. Hover: `rosewood-muted` background.

**Cards / Containers**
- Main cards (student profiles, history): `rounded-xl` (1.5rem), Ivory background, 1px border, whisper-soft shadow `0 1px 4px rgba(28,25,23,0.06)`.
- Standard UI (buttons, inputs, list items): `rounded` (0.5rem).
- Badges (level indicators): `0.375rem` radius — compact feel.
- High-density lists: no card background — `border-top` dividers, `rosewood-muted` on row hover.

**Inputs / Forms**
- Labels always positioned above the input.
- Focus: 2px solid Rosewood ring, 2px offset.
- Helper text: Graphite. Error state: Destructive Clay for both text and border.

**Assessment Slider** (core interaction component)
- Track: full width. Active portion solid Rosewood; inactive portion Whisper Border.
- Thumb: 28px circle, Ivory fill + 2px Rosewood border by default. Dragging: scales to 34px, switches to solid Rosewood fill.
- Feedback: score numerals (Geist Mono) animate/counter in real-time as slider moves.
- Minimum 44px vertical touch zone (one-handed mobile operation).

**Feedback & Status**
- Level badges: Bronze/Silver/Gold at 15% opacity background with matching text color.
- Offline banner: fixed 44px bottom strip, Caution Amber. Success: Success Sage, auto-dismiss.
- Skeleton loaders: shimmer matching target component dimensions — no circular spinners.

## 5. Layout Principles

Fluid Grid model centered on the user's focus — the student assessment.

- **Grid System:** CSS Grid throughout. Content constrained to `1400px` max-width.
- **Asymmetry:** Favor `2-column` or `2fr 1fr` splits over equal-width columns to guide the eye toward critical data.
- **Mobile First:** Below `768px`, all layouts collapse to single column. Assessment Panel uses stacked full-width sliders for one-handed operation.
- **Touch Targets:** All interactive elements — especially sliders — maintain minimum `44px` vertical touch zone.
- **Tables:** Desktop uses traditional rows with `1px` top dividers; mobile transforms rows into stacked card-recap components.
- **Spacing scale:** section-gap `clamp(3rem, 6vw, 5rem)`, gutter `1rem`, margin-mobile `1rem`, margin-desktop `2.5rem`.

## 6. Elevation & Depth

Depth indicates functional hierarchy, not aesthetic flourish.
- Primary mechanism: tonal contrast between `parchment` background and `ivory-surface` containers.
- Ambient shadows only on primary cards: `0 1px 4px rgba(28,25,23,0.06)`.
- Most functional elements (inputs, secondary buttons) rely on `1px whisper-border` outlines, not shadows.
- No outer glows, neon focus effects, or heavy drop shadows.

## 7. Shapes

Shape language is **Rounded** — a soft counterpoint to typographic precision.
- Main cards: `rounded-xl` (1.5rem).
- Standard UI: `rounded` (0.5rem).
- Badges: `0.375rem`.
- Slider thumbs: perfect circles.
