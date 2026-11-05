import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../data/repositories/auth_repository.dart';
import 'providers/auth_provider.dart';
import 'widgets/linkedin_auth_button.dart';

class AuthCallbackScreen extends ConsumerStatefulWidget {
  const AuthCallbackScreen({
    this.code,
    this.error,
    super.key,
  });

  final String? code;
  final String? error;

  @override
  ConsumerState<AuthCallbackScreen> createState() => _AuthCallbackScreenState();
}

class _AuthCallbackScreenState extends ConsumerState<AuthCallbackScreen> {
  String? _errorMessage;
  bool _isProcessing = true;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _processAuthCode());
  }

  Future<void> _processAuthCode() async {
    if (widget.error != null && widget.error!.isNotEmpty) {
      setState(() {
        _errorMessage = 'LinkedIn authentication was cancelled or rejected.';
        _isProcessing = false;
      });
      return;
    }

    final code = widget.code;
    if (code == null || code.isEmpty) {
      setState(() {
        _errorMessage = 'No authorization code received from LinkedIn.';
        _isProcessing = false;
      });
      return;
    }

    try {
      var redirectUri = linkedInRedirectUri;
      try {
        final config = await ref.read(authRepositoryProvider).getAuthConfig();
        if (config.redirectUri.isNotEmpty) {
          redirectUri = config.redirectUri;
        }
      } catch (_) {
        // Fall back to default if backend config unreachable
      }

      await ref.read(authNotifierProvider.notifier).loginWithLinkedInCode(
            code: code,
            redirectUri: redirectUri,
          );
      // Navigation is triggered automatically by router redirect based on auth state
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Failed to verify account: $e';
          _isProcessing = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: _isProcessing
              ? Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const CircularProgressIndicator(color: AppColors.cinnamon),
                    const SizedBox(height: 24),
                    Text(
                      'Verifying with LinkedIn...',
                      style: theme.textTheme.titleMedium?.copyWith(
                        color: isDark ? AppColors.cream : AppColors.ink,
                      ),
                    ),
                  ],
                )
              : Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.error_outline_rounded,
                      size: 56,
                      color: isDark ? AppColors.errorDark : AppColors.error,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      'Authentication Failed',
                      style: theme.textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: isDark ? AppColors.cream : AppColors.ink,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      _errorMessage ?? 'An error occurred during authentication.',
                      textAlign: TextAlign.center,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: isDark
                            ? AppColors.cream.withAlpha(180)
                            : AppColors.ink.withAlpha(180),
                      ),
                    ),
                    const SizedBox(height: 24),
                    ElevatedButton(
                      onPressed: () => context.go('/welcome'),
                      style: ElevatedButton.styleFrom(
                        backgroundColor:
                            isDark ? AppColors.cinnamon : AppColors.brown,
                        foregroundColor:
                            isDark ? AppColors.ink : AppColors.cream,
                      ),
                      child: const Text('Back to Welcome'),
                    ),
                  ],
                ),
        ),
      ),
    );
  }
}
