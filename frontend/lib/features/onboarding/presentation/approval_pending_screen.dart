import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../auth/presentation/providers/auth_provider.dart';

class ApprovalPendingScreen extends ConsumerWidget {
  const ApprovalPendingScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Spacer(),
              Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  color: isDark
                      ? AppColors.espressoSurface
                      : AppColors.cinnamon.withAlpha(30),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.hourglass_top_rounded,
                  size: 48,
                  color: AppColors.cinnamon,
                ),
              ),
              const SizedBox(height: 24),
              Text(
                'Approval Pending',
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: isDark ? AppColors.cream : AppColors.ink,
                ),
              ),
              const SizedBox(height: 12),
              Text(
                'Your account has been created and is waiting for verification by the Alumni Association.',
                textAlign: TextAlign.center,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: isDark
                      ? AppColors.cream.withAlpha(180)
                      : AppColors.ink.withAlpha(180),
                  height: 1.4,
                ),
              ),
              const Spacer(),
              TextButton.icon(
                onPressed: () =>
                    ref.read(authNotifierProvider.notifier).logout(),
                icon: Icon(
                  Icons.logout_rounded,
                  color: isDark ? AppColors.errorDark : AppColors.error,
                ),
                label: Text(
                  'Log Out',
                  style: TextStyle(
                    color: isDark ? AppColors.errorDark : AppColors.error,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
