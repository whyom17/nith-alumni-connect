import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../app/app_drawer.dart';
import '../../../../core/network/api_client.dart';
import '../../../../core/network/health_repository.dart';

/// The Home screen for NITH Alumni Connect, featuring the announcements feed
/// and the M0 backend health check status card.
class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final healthState = ref.watch(healthCheckProvider);

    return Scaffold(
      drawer: const AppDrawer(),
      appBar: AppBar(
        title: const Text('Alumni Connect'),
        leading: Builder(
          builder: (scaffoldContext) {
            return IconButton(
              icon: const Icon(Icons.menu_rounded),
              tooltip: 'Open navigation menu',
              onPressed: () => Scaffold.of(scaffoldContext).openDrawer(),
            );
          },
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.inbox_outlined),
            tooltip: 'Inbox',
            onPressed: () => context.push('/inbox'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.refresh(healthCheckProvider.future),
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Backend Connectivity Diagnostic Card (M0 Foundation)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.network_check_rounded,
                          color: colorScheme.primary,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Backend Connection (M0)',
                          style: theme.textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const Spacer(),
                        IconButton(
                          icon: const Icon(Icons.refresh_rounded, size: 20),
                          tooltip: 'Check backend health',
                          onPressed: () => ref.invalidate(healthCheckProvider),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'API Base URL: $defaultApiBaseUrl',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: colorScheme.secondary,
                      ),
                    ),
                    const SizedBox(height: 12),
                    healthState.when(
                      data: (status) => Row(
                        children: [
                          const Icon(
                            Icons.check_circle_rounded,
                            color: Colors.green,
                            size: 20,
                          ),
                          const SizedBox(width: 8),
                          Text(
                            'Status: $status (Connected to /healthz)',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: Colors.green.shade800,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                      loading: () => const Row(
                        children: [
                          SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          ),
                          SizedBox(width: 12),
                          Text('Checking /healthz...'),
                        ],
                      ),
                      error: (err, _) => Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Icon(
                            Icons.error_outline_rounded,
                            color: colorScheme.error,
                            size: 20,
                          ),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              'Failed to connect: $err',
                              style: TextStyle(color: colorScheme.error),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            // Home Feed Announcements Placeholder
            Center(
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 32),
                child: Column(
                  children: [
                    Icon(
                      Icons.newspaper_rounded,
                      size: 48,
                      color: colorScheme.secondary,
                    ),
                    const SizedBox(height: 12),
                    Text(
                      'Announcements Feed',
                      style: theme.textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Posts, campus updates, jobs, and events will appear here.',
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: colorScheme.secondary,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
