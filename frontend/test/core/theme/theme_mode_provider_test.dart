import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:nith_alumni_connect/core/theme/theme_mode_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  group('ThemeModeNotifier', () {
    test('defaults to ThemeMode.system when no preference is saved', () async {
      SharedPreferences.setMockInitialValues({});
      final prefs = await SharedPreferences.getInstance();

      final container = ProviderContainer(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
        ],
      );
      addTearDown(container.dispose);

      final themeMode = container.read(themeModeProvider);
      expect(themeMode, ThemeMode.system);
    });

    test('loads saved dark mode preference', () async {
      SharedPreferences.setMockInitialValues({'nith_theme_mode': 'dark'});
      final prefs = await SharedPreferences.getInstance();

      final container = ProviderContainer(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
        ],
      );
      addTearDown(container.dispose);

      final themeMode = container.read(themeModeProvider);
      expect(themeMode, ThemeMode.dark);
    });

    test('setThemeMode updates state and persists value', () async {
      SharedPreferences.setMockInitialValues({});
      final prefs = await SharedPreferences.getInstance();

      final container = ProviderContainer(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
        ],
      );
      addTearDown(container.dispose);

      final notifier = container.read(themeModeProvider.notifier);
      await notifier.setThemeMode(ThemeMode.light);

      expect(container.read(themeModeProvider), ThemeMode.light);
      expect(prefs.getString('nith_theme_mode'), 'light');
    });

    test('toggle switches from light to dark', () async {
      SharedPreferences.setMockInitialValues({'nith_theme_mode': 'light'});
      final prefs = await SharedPreferences.getInstance();

      final container = ProviderContainer(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
        ],
      );
      addTearDown(container.dispose);

      final notifier = container.read(themeModeProvider.notifier);
      await notifier.toggle(Brightness.light);

      expect(container.read(themeModeProvider), ThemeMode.dark);
      expect(prefs.getString('nith_theme_mode'), 'dark');
    });
  });
}
