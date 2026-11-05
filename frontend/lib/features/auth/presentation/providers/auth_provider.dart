import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/storage/token_storage.dart';
import '../../data/models/auth_dto.dart';
import '../../data/repositories/auth_repository.dart';
import '../../domain/models/auth_state.dart';

class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier({
    required this.authRepository,
    required this.tokenStorage,
  }) : super(const AuthState.initial()) {
    initSession();
  }

  final AuthRepository authRepository;
  final TokenStorage tokenStorage;

  /// Restores active session on app startup using stored refresh token.
  Future<void> initSession() async {
    try {
      final refreshToken = await tokenStorage.getRefreshToken();
      if (refreshToken == null || refreshToken.isEmpty) {
        state = const AuthState.unauthenticated();
        return;
      }

      final resp = await authRepository.refreshToken(refreshToken);
      await tokenStorage.saveTokens(
        accessToken: resp.accessToken,
        refreshToken: resp.refreshToken,
      );

      // Session restored: create basic authenticated state
      // (Full profile will be loaded by meEndpoint in M2)
      state = const AuthState.authenticated(
        AuthUserDto(
          id: '',
          email: '',
          role: 'alumnus',
          status: 'pending',
          profileComplete: false,
        ),
      );
    } catch (_) {
      await tokenStorage.clearTokens();
      state = const AuthState.unauthenticated();
    }
  }

  /// Exchanges LinkedIn code for tokens and persists session.
  Future<AuthResponseDto> loginWithLinkedInCode({
    required String code,
    required String redirectUri,
    String? codeVerifier,
  }) async {
    state = const AuthState.loading();
    try {
      final resp = await authRepository.authenticateWithLinkedIn(
        code: code,
        redirectUri: redirectUri,
        codeVerifier: codeVerifier,
      );

      await tokenStorage.saveTokens(
        accessToken: resp.accessToken,
        refreshToken: resp.refreshToken,
      );

      state = AuthState.authenticated(resp.user);
      return resp;
    } catch (e) {
      state = AuthState.unauthenticated(e.toString());
      rethrow;
    }
  }

  /// Logs out of session and removes stored tokens.
  Future<void> logout() async {
    try {
      final refreshToken = await tokenStorage.getRefreshToken();
      if (refreshToken != null && refreshToken.isNotEmpty) {
        await authRepository.logout(refreshToken);
      }
    } catch (_) {
      // Best-effort logout; proceed to clear local state
    } finally {
      await tokenStorage.clearTokens();
      state = const AuthState.unauthenticated();
    }
  }
}

final authNotifierProvider =
    StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final authRepo = ref.watch(authRepositoryProvider);
  final tokenStorage = ref.watch(tokenStorageProvider);
  return AuthNotifier(
    authRepository: authRepo,
    tokenStorage: tokenStorage,
  );
});
