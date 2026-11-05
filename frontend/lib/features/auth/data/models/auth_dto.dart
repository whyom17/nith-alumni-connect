/// Public configuration returned from GET /v1/auth/config.
class AuthConfigDto {
  const AuthConfigDto({
    required this.clientId,
    required this.redirectUri,
  });

  factory AuthConfigDto.fromJson(Map<String, dynamic> json) {
    return AuthConfigDto(
      clientId: json['client_id'] as String? ?? '',
      redirectUri: json['redirect_uri'] as String? ?? '',
    );
  }

  final String clientId;
  final String redirectUri;
}

/// Data transfer object for the authenticated user payload.
class AuthUserDto {
  const AuthUserDto({
    required this.id,
    required this.email,
    required this.role,
    required this.status,
    required this.profileComplete,
  });

  factory AuthUserDto.fromJson(Map<String, dynamic> json) {
    return AuthUserDto(
      id: json['id'] as String,
      email: json['email'] as String,
      role: json['role'] as String,
      status: json['status'] as String,
      profileComplete: json['profile_complete'] as bool? ?? false,
    );
  }

  final String id;
  final String email;
  final String role;
  final String status;
  final bool profileComplete;

  Map<String, dynamic> toJson() => {
        'id': id,
        'email': email,
        'role': role,
        'status': status,
        'profile_complete': profileComplete,
      };
}

/// Response returned from POST /v1/auth/linkedin.
class AuthResponseDto {
  const AuthResponseDto({
    required this.accessToken,
    required this.refreshToken,
    required this.isNewUser,
    required this.user,
  });

  factory AuthResponseDto.fromJson(Map<String, dynamic> json) {
    return AuthResponseDto(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      isNewUser: json['is_new_user'] as bool? ?? false,
      user: AuthUserDto.fromJson(json['user'] as Map<String, dynamic>),
    );
  }

  final String accessToken;
  final String refreshToken;
  final bool isNewUser;
  final AuthUserDto user;
}

/// Response returned from POST /v1/auth/refresh.
class RefreshTokenResponseDto {
  const RefreshTokenResponseDto({
    required this.accessToken,
    required this.refreshToken,
  });

  factory RefreshTokenResponseDto.fromJson(Map<String, dynamic> json) {
    return RefreshTokenResponseDto(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
    );
  }

  final String accessToken;
  final String refreshToken;
}
