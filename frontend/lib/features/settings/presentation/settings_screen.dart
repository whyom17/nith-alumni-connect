import 'package:flutter/material.dart';

/// Settings & Privacy screen (Phase 1 M7 placeholder).
class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
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
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Log Out will be wired in M1.')),
              );
            },
          ),
        ],
      ),
    );
  }
}
