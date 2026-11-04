import 'package:flutter/material.dart';

/// The Search Alumni tab screen (Phase 1 M3 placeholder).
class SearchScreen extends StatelessWidget {
  const SearchScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Search Alumni'),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.search_rounded,
                size: 56,
                color: colorScheme.secondary,
              ),
              const SizedBox(height: 16),
              Text(
                'Alumni Directory',
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Find alumni by name, batch, branch, city, and industry.\n(Coming in Milestone 3)',
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
