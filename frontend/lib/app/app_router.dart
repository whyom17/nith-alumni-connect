import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../features/announcements/presentation/home_screen.dart';
import '../features/announcements/presentation/saved_posts_screen.dart';
import '../features/directory/presentation/search_screen.dart';
import '../features/groups/presentation/groups_screen.dart';
import '../features/inbox/presentation/inbox_screen.dart';
import '../features/profile/presentation/profile_screen.dart';
import '../features/settings/presentation/settings_screen.dart';
import 'app_shell.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>(debugLabel: 'root');

/// The central GoRouter provider configured with the 4-tab StatefulShellRoute.
final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/home',
    debugLogDiagnostics: false,
    routes: [
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) {
          return AppShell(navigationShell: navigationShell);
        },
        branches: [
          // Tab 1: Home (Announcements feed)
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/home',
                builder: (context, state) => const HomeScreen(),
              ),
            ],
          ),
          // Tab 2: Groups (Phase 1 Coming Soon)
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/groups',
                builder: (context, state) => const GroupsScreen(),
              ),
            ],
          ),
          // Tab 3: Search Alumni
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/search',
                builder: (context, state) => const SearchScreen(),
              ),
            ],
          ),
          // Tab 4: My Profile
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/profile',
                builder: (context, state) => const ProfileScreen(),
              ),
            ],
          ),
        ],
      ),
      // Standalone full-screen routes pushed above bottom navigation
      GoRoute(
        path: '/inbox',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const InboxScreen(),
      ),
      GoRoute(
        path: '/saved',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SavedPostsScreen(),
      ),
      GoRoute(
        path: '/settings',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SettingsScreen(),
      ),
    ],
  );
});
