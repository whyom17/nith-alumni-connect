import 'package:flutter/material.dart';

/// The Saved Posts screen accessible from the drawer (Phase 1 M5).
class SavedPostsScreen extends StatelessWidget {
  const SavedPostsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Saved Posts'),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.bookmark_border_rounded,
                size: 56,
                color: colorScheme.secondary,
              ),
              const SizedBox(height: 16),
              Text(
                'Saved Announcements',
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Bookmarks of jobs, events, and initiatives will appear here.\n(Coming in Milestone 5)',
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
