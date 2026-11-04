import 'package:flutter/material.dart';

/// The Inbox screen for connection requests and notifications (Phase 1 M6).
class InboxScreen extends StatelessWidget {
  const InboxScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Inbox'),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.notifications_none_rounded,
                size: 56,
                color: colorScheme.secondary,
              ),
              const SizedBox(height: 16),
              Text(
                'Notifications & Requests',
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Connection requests and activity alerts will appear here.\n(Coming in Milestone 6)',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: colorScheme.secondary,
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
