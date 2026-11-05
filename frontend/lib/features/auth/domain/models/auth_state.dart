import '../../data/models/auth_dto.dart';

enum AuthStatus {
  initial,
  loading,
  authenticated,
  unauthenticated,
}

class AuthState {
  const AuthState({
    required this.status,
    this.user,
    this.errorMessage,
  });

  const AuthState.initial()
      : status = AuthStatus.initial,
        user = null,
        errorMessage = null;

  const AuthState.unauthenticated([this.errorMessage])
      : status = AuthStatus.unauthenticated,
        user = null;

  const AuthState.loading()
      : status = AuthStatus.loading,
        user = null,
        errorMessage = null;

  const AuthState.authenticated(this.user)
      : status = AuthStatus.authenticated,
        errorMessage = null;

  final AuthStatus status;
  final AuthUserDto? user;
  final String? errorMessage;

  bool get isAuthenticated => status == AuthStatus.authenticated && user != null;
  bool get isApproved => user?.status == 'approved';
  bool get isPending => user?.status == 'pending';
  bool get isInitialLoading => status == AuthStatus.initial;
}
