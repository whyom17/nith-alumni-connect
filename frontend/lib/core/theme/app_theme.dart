import 'package:flutter/material.dart';
import 'app_colors.dart';

/// Provides light and dark ThemeData configured for the Spiced Chai design system.
abstract final class AppTheme {
  static ThemeData get lightTheme {
    final colorScheme = ColorScheme(
      brightness: Brightness.light,
      primary: AppColors.brown,
      onPrimary: AppColors.cream,
      secondary: AppColors.nutmeg,
      onSecondary: AppColors.cream,
      tertiary: AppColors.cinnamon,
      onTertiary: AppColors.ink,
      error: AppColors.error,
      onError: Colors.white,
      surface: AppColors.surface,
      onSurface: AppColors.ink,
      outline: AppColors.nutmeg.withAlpha(102), // 40% opacity
    );

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: AppColors.cream,
      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.cream,
        foregroundColor: AppColors.ink,
        elevation: 0,
        centerTitle: false,
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: AppColors.surface,
        indicatorColor: AppColors.cinnamon.withAlpha(51), // 20% opacity
        iconTheme: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const IconThemeData(color: AppColors.brown);
          }
          return const IconThemeData(color: AppColors.nutmeg);
        }),
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const TextStyle(
              color: AppColors.brown,
              fontWeight: FontWeight.w600,
            );
          }
          return const TextStyle(color: AppColors.nutmeg);
        }),
      ),
      drawerTheme: const DrawerThemeData(
        backgroundColor: AppColors.surface,
      ),
      cardTheme: CardThemeData(
        color: AppColors.surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: BorderSide(
            color: AppColors.nutmeg.withAlpha(51), // 20% opacity
            width: 1,
          ),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: AppColors.brown,
          foregroundColor: AppColors.cream,
          minimumSize: const Size.fromHeight(48),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: AppColors.surface,
        selectedColor: AppColors.cinnamon,
        labelStyle: const TextStyle(color: AppColors.ink),
        secondaryLabelStyle: const TextStyle(color: AppColors.ink),
        side: const BorderSide(color: AppColors.nutmeg, width: 1),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
      dividerTheme: DividerThemeData(
        color: AppColors.nutmeg.withAlpha(51),
        thickness: 1,
      ),
    );
  }

  static ThemeData get darkTheme {
    final colorScheme = ColorScheme(
      brightness: Brightness.dark,
      primary: AppColors.cinnamon,
      onPrimary: AppColors.ink,
      secondary: AppColors.nutmeg,
      onSecondary: AppColors.cream,
      tertiary: AppColors.cinnamon,
      onTertiary: AppColors.ink,
      error: AppColors.errorDark,
      onError: AppColors.ink,
      surface: AppColors.espressoSurface,
      onSurface: AppColors.cream,
      outline: AppColors.nutmeg.withAlpha(153), // 60% opacity
    );

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: AppColors.espresso,
      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.espresso,
        foregroundColor: AppColors.cream,
        elevation: 0,
        centerTitle: false,
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: AppColors.espressoSurface,
        indicatorColor: AppColors.cinnamon.withAlpha(51), // 20% opacity
        iconTheme: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const IconThemeData(color: AppColors.cinnamon);
          }
          return const IconThemeData(color: AppColors.nutmeg);
        }),
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return const TextStyle(
              color: AppColors.cinnamon,
              fontWeight: FontWeight.w600,
            );
          }
          return const TextStyle(color: AppColors.nutmeg);
        }),
      ),
      drawerTheme: const DrawerThemeData(
        backgroundColor: AppColors.espressoSurface,
      ),
      cardTheme: CardThemeData(
        color: AppColors.espressoSurface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: BorderSide(
            color: AppColors.nutmeg.withAlpha(51),
            width: 1,
          ),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: AppColors.cinnamon,
          foregroundColor: AppColors.ink,
          minimumSize: const Size.fromHeight(48),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
      ),
      chipTheme: ChipThemeData(
        backgroundColor: AppColors.espressoSurface,
        selectedColor: AppColors.cinnamon,
        labelStyle: const TextStyle(color: AppColors.cream),
        secondaryLabelStyle: const TextStyle(color: AppColors.ink),
        side: const BorderSide(color: AppColors.nutmeg, width: 1),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
      dividerTheme: DividerThemeData(
        color: AppColors.nutmeg.withAlpha(51),
        thickness: 1,
      ),
    );
  }
}
