import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../data/repositories/auth_repository.dart';
import '../providers/auth_provider.dart';

const String linkedInClientId = String.fromEnvironment(
  'LINKEDIN_CLIENT_ID',
  defaultValue: '',
);

const String linkedInRedirectUri = String.fromEnvironment(
  'LINKEDIN_REDIRECT_URI',
  defaultValue: 'org.nith.alumni://callback',
);

class LinkedInAuthButton extends ConsumerStatefulWidget {
  const LinkedInAuthButton({
    required this.label,
    super.key,
  });

  final String label;

  @override
  ConsumerState<LinkedInAuthButton> createState() => _LinkedInAuthButtonState();
}

class _LinkedInAuthButtonState extends ConsumerState<LinkedInAuthButton> {
  bool _isLoading = false;

  Future<void> _handleLinkedInAuth() async {
    setState(() => _isLoading = true);
    String? connectionError;

    try {
      var clientId = linkedInClientId;
      var redirectUri = linkedInRedirectUri;

      // If client ID was not supplied via compile-time --dart-define, fetch it from backend
      if (clientId.isEmpty) {
        try {
          final authConfig =
              await ref.read(authRepositoryProvider).getAuthConfig();
          clientId = authConfig.clientId;
          if (authConfig.redirectUri.isNotEmpty) {
            redirectUri = authConfig.redirectUri;
          }
        } catch (e) {
          connectionError = e.toString();
        }
      }

      if (clientId.isEmpty) {
        if (connectionError != null && mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                'Could not reach backend to load LinkedIn Client ID ($connectionError). Please ensure the backend is running.',
              ),
              duration: const Duration(seconds: 4),
            ),
          );
        }
        _showCodeCompletionDialog(redirectUri, errorDetail: connectionError);
        return;
      }

      final authUrl = Uri.https('www.linkedin.com', '/oauth/v2/authorization', {
        'response_type': 'code',
        'client_id': clientId,
        'redirect_uri': redirectUri,
        'scope': 'openid profile email',
        'state': 'nith_auth_${DateTime.now().millisecondsSinceEpoch}',
      });

      final launched = await launchUrl(
        authUrl,
        mode: LaunchMode.externalApplication,
      );

      if (launched && mounted) {
        _showCodeCompletionDialog(redirectUri);
      } else if (!launched && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Could not open LinkedIn login URL.')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to start login: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showCodeCompletionDialog(String redirectUri, {String? errorDetail}) {
    final controller = TextEditingController();
    showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('Complete Sign In'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              errorDetail != null
                  ? 'Unable to connect to backend automatically ($errorDetail).'
                  : 'After authenticating in your browser, if the app did not open automatically, copy the authorization code from the browser and paste it here:',
              style: const TextStyle(fontSize: 13),
            ),
            const SizedBox(height: 14),
            TextField(
              controller: controller,
              decoration: const InputDecoration(
                hintText: 'Paste code here...',
                border: OutlineInputBorder(),
                prefixIcon: Icon(Icons.vpn_key_rounded),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogCtx),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              final code = controller.text.trim();
              if (code.isEmpty) return;
              Navigator.pop(dialogCtx);
              try {
                await ref
                    .read(authNotifierProvider.notifier)
                    .loginWithLinkedInCode(
                      code: code,
                      redirectUri: redirectUri,
                    );
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Authentication failed: $e')),
                  );
                }
              }
            },
            child: const Text('Submit'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: _isLoading ? null : _handleLinkedInAuth,
      style: ElevatedButton.styleFrom(
        backgroundColor: const Color(0xFF0A66C2), // Official LinkedIn blue
        foregroundColor: Colors.white,
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        elevation: 0,
      ),
      child: _isLoading
          ? const SizedBox(
              height: 20,
              width: 20,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: Colors.white,
              ),
            )
          : Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  padding: const EdgeInsets.all(2),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: const Text(
                    ' in ',
                    style: TextStyle(
                      color: Color(0xFF0A66C2),
                      fontWeight: FontWeight.w900,
                      fontSize: 13,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Text(
                  widget.label,
                  style: const TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
    );
  }
}
