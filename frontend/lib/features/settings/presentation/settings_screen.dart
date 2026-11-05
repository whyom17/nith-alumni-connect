import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../auth/presentation/providers/auth_provider.dart';

/// Settings & Privacy screen.
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings & Privacy'),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Text(
            'Privacy & Communication',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.primary,
            ),
          ),
          const SizedBox(height: 8),
          SwitchListTile(
            title: const Text('Allow Direct WhatsApp Contact'),
            subtitle: const Text('When connected, show WhatsApp button on profile'),
            value: true,
            onChanged: (val) {},
          ),
          SwitchListTile(
            title: const Text('Contact Details Visibility'),
            subtitle: const Text('Show phone and email to accepted connections'),
            value: false,
            onChanged: (val) {},
          ),
          const Divider(height: 32),
          Text(
            'About',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.primary,
            ),
          ),
          const ListTile(
            title: Text('NIT Hamirpur Alumni Connect'),
            subtitle: Text('Version 1.0.0 (Phase 1 MVP)'),
          ),
          const Divider(height: 32),
          ListTile(
            leading: Icon(Icons.logout_rounded, color: colorScheme.error),
            title: Text(
              'Log Out',
              style: TextStyle(
                color: colorScheme.error,
                fontWeight: FontWeight.bold,
              ),
            ),
            onTap: () async {
              await ref.read(authNotifierProvider.notifier).logout();
            },
          ),
        ],
      ),
    );
  }
}
