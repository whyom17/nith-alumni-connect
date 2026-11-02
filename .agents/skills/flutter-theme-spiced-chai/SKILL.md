---
name: flutter-theme-spiced-chai
description: The Spiced Chai visual design system (cream #FDFBD4, cinnamon #D47E30, brown #8D5A2B, nutmeg #825E34) for the Alumni Connect Flutter app, in both light and dark modes with a user toggle. Use whenever writing or editing any UI, widget, screen, theme, dark mode, color, text style, button, card, chip, badge, icon, or app bar, or when asked about styling, branding, contrast, or look and feel.
---

# Spiced Chai Theme (light + dark)

Warm, cozy, approachable. The app ships **light and dark**; the user switches via the circle toggle in the drawer. Default = follow system. Persist only the mode (`shared_preferences`).

## Tokens (`lib/core/theme/app_colors.dart`)
```dart
abstract final class AppColors {
  // Brand (from the palette)
  static const cream = Color(0xFFFDFBD4);
  static const cinnamon = Color(0xFFD47E30);
  static const brown = Color(0xFF8D5A2B);
  static const nutmeg = Color(0xFF825E34);
  // Derived
  static const ink = Color(0xFF3B2712);
  static const surface = Color(0xFFFFFDEA);
  static const espresso = Color(0xFF1C140D);
  static const espressoSurface = Color(0xFF2A1F15);
  static const error = Color(0xFFB3261E);
  static const errorDark = Color(0xFFFF8A80);
}
```

## ColorScheme mapping
| Role | Light | Dark |
|---|---|---|
| scaffold background | cream | espresso |
| surface | surface | espressoSurface |
| onSurface | ink | cream |
| primary / onPrimary | brown / cream | cinnamon / ink |
| secondary / onSecondary | nutmeg / cream | nutmeg / cream |
| tertiary (accent) / onTertiary | cinnamon / ink | cinnamon / ink |
| outline | nutmeg @ 40% | nutmeg @ 60% |
| error | error | errorDark |

## Component rules
- **Filled button (Connect, Get Started, Sign up with LinkedIn):** primary, 12dp radius, 48dp min height. In Search rows use a compact variant. Connected/Pending use a tonal (outlined) style.
- **Contact icon buttons (WhatsApp, Email):** circular, tonal, with clear icons. Do not use brand green/yellow.
- **Bottom `NavigationBar`:** surface background, selected indicator cinnamon @ 20%, selected icon/label primary.
- **Drawer:** surface background; the dark/light toggle is an icon button (sun/moon) top-right of the header.
- **Post card:** surface, 16dp radius, 1dp nutmeg@20% border; dividers between header/body/actions.
- **Chips (filters):** unselected outlined nutmeg; selected cinnamon fill with ink text.
- **Verified badge:** diamond icon in cinnamon (thin brown outline in light mode), `Semantics(label: 'Verified alumnus')`.
- **Log Out:** error colour text.
- **Text fields:** surface fill, nutmeg border, primary focus border.
- Spacing scale 4, 8, 12, 16, 24, 32.

## Accessibility
- Body text at least 4.5:1. White/cream on cinnamon fails for small text, so use ink on cinnamon.
- brown and nutmeg are very close; don't rely on the difference.
- Tap targets 48dp; text scale up to 1.3; verify both modes.

## Implementation
Build `lightTheme` and `darkTheme` in `lib/core/theme/app_theme.dart` (`useMaterial3: true`) with component themes (`FilledButtonTheme`, `CardTheme`, `InputDecorationTheme`, `ChipTheme`, `AppBarTheme`, `NavigationBarTheme`, `DrawerTheme`, `DividerTheme`). Wire `MaterialApp.router(theme:, darkTheme:, themeMode:)`. Widgets never hardcode hex colours.
