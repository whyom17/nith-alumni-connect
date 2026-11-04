import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// Provider for SharedPreferences instance, initialized in main.dart bootstrap.
final sharedPreferencesProvider = Provider<SharedPreferences>((ref) {
  throw UnimplementedError(
    'sharedPreferencesProvider must be overridden in ProviderScope',
  );
});

/// Controls and persists the application's [ThemeMode] (System, Light, or Dark).
class ThemeModeNotifier extends Notifier<ThemeMode> {
  static const _storageKey = 'nith_theme_mode';

  @override
  ThemeMode build() {
    final prefs = ref.watch(sharedPreferencesProvider);
    final savedValue = prefs.getString(_storageKey);

    return switch (savedValue) {
      'light' => ThemeMode.light,
      'dark' => ThemeMode.dark,
      _ => ThemeMode.system,
    };
  }

  /// Sets an explicit [ThemeMode] and persists the choice.
  Future<void> setThemeMode(ThemeMode mode) async {
    state = mode;
    final prefs = ref.read(sharedPreferencesProvider);

    switch (mode) {
      case ThemeMode.light:
        await prefs.setString(_storageKey, 'light');
      case ThemeMode.dark:
        await prefs.setString(_storageKey, 'dark');
      case ThemeMode.system:
        await prefs.remove(_storageKey);
    }
  }

  /// Toggles between light and dark modes based on the current active brightness.
  Future<void> toggle(Brightness currentPlatformBrightness) async {
    final isCurrentlyDark = switch (state) {
      ThemeMode.dark => true,
      ThemeMode.light => false,
      ThemeMode.system => currentPlatformBrightness == Brightness.dark,
    };

    await setThemeMode(isCurrentlyDark ? ThemeMode.light : ThemeMode.dark);
  }
}

final themeModeProvider =
    NotifierProvider<ThemeModeNotifier, ThemeMode>(ThemeModeNotifier.new);
