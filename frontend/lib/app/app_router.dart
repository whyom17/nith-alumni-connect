import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../features/announcements/presentation/home_screen.dart';
import '../features/announcements/presentation/saved_posts_screen.dart';
import '../features/auth/presentation/auth_callback_screen.dart';
import '../features/auth/presentation/create_account_screen.dart';
import '../features/auth/presentation/login_screen.dart';
import '../features/auth/presentation/providers/auth_provider.dart';
import '../features/auth/presentation/welcome_screen.dart';
import '../features/directory/presentation/search_screen.dart';
import '../features/groups/presentation/groups_screen.dart';
import '../features/inbox/presentation/inbox_screen.dart';
import '../features/onboarding/presentation/approval_pending_screen.dart';
import '../features/profile/presentation/profile_screen.dart';
import '../features/settings/presentation/settings_screen.dart';
import 'app_shell.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>(debugLabel: 'root');

/// The central GoRouter provider configured with auth guards and tab navigation.
final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authNotifierProvider);

  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/welcome',
    debugLogDiagnostics: false,
    redirect: (context, state) {
      // While checking stored credentials, don't redirect
      if (authState.isInitialLoading) {
        return null;
      }

      final loc = state.matchedLocation;
      final isAuthFlow = loc == '/welcome' ||
          loc == '/signup' ||
          loc == '/login' ||
          loc == '/callback';

      if (!authState.isAuthenticated) {
        // If unauthenticated and on a protected screen, redirect to /welcome
        return isAuthFlow ? null : '/welcome';
      }

      // If authenticated and visiting auth screens, redirect to /home
      if (isAuthFlow) {
        return '/home';
      }

      return null;
    },
    routes: [
      // Auth screens
      GoRoute(
        path: '/welcome',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const WelcomeScreen(),
      ),
      GoRoute(
        path: '/signup',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const CreateAccountScreen(),
      ),
      GoRoute(
        path: '/login',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/callback',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) {
          final code = state.uri.queryParameters['code'];
          final error = state.uri.queryParameters['error'];
          return AuthCallbackScreen(code: code, error: error);
        },
      ),
      GoRoute(
        path: '/approval-pending',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ApprovalPendingScreen(),
      ),

      // Main app bottom navigation shell
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
