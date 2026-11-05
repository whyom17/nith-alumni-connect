import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../storage/token_storage.dart';
import 'app_failure.dart';
import 'auth_interceptor.dart';

/// Base URL passed at compile time via `--dart-define=API_BASE_URL=...`.
/// Defaults to Android emulator loopback `http://10.0.2.2:8080` on Android,
/// or `http://localhost:8080` on web, Windows, macOS, and iOS simulator.
const String _envBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: '');

String get defaultApiBaseUrl {
  if (_envBaseUrl.isNotEmpty) return _envBaseUrl;
  if (kIsWeb) return 'http://localhost:8080';
  if (defaultTargetPlatform == TargetPlatform.android) {
    return 'http://10.0.2.2:8080';
  }
  return 'http://localhost:8080';
}

/// Maps Dio errors into typed domain [AppFailure] instances.
AppFailure mapDioException(DioException e) {
  if (e.type == DioExceptionType.connectionTimeout ||
      e.type == DioExceptionType.sendTimeout ||
      e.type == DioExceptionType.receiveTimeout ||
      e.type == DioExceptionType.connectionError) {
    return const NetworkFailure(
      'Unable to reach server. Please check your connection.',
    );
  }

  final response = e.response;
  if (response != null && response.data is Map<String, dynamic>) {
    final data = response.data as Map<String, dynamic>;
    if (data['error'] is Map<String, dynamic>) {
      final err = data['error'] as Map<String, dynamic>;
      final message = err['message'] as String? ?? 'Server returned an error';
      final fieldsRaw = err['fields'];
      final fields = fieldsRaw is Map<String, dynamic>
          ? fieldsRaw.map((k, v) => MapEntry(k, v.toString()))
          : null;

      return switch (response.statusCode) {
        400 => ValidationFailure(message, fields: fields),
        401 => UnauthorizedFailure(message),
        403 => ForbiddenFailure(message),
        404 => NotFoundFailure(message),
        500 || 502 || 503 => ServerFailure(message),
        _ => UnknownFailure(message),
      };
    }
  }

  return UnknownFailure(e.message ?? 'An unexpected error occurred');
}

/// Provides a preconfigured [Dio] client for backend API communication.
final dioProvider = Provider<Dio>((ref) {
  final tokenStorage = ref.watch(tokenStorageProvider);
  final dio = Dio(
    BaseOptions(
      baseUrl: defaultApiBaseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      sendTimeout: const Duration(seconds: 10),
      headers: const {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  dio.interceptors.add(
    AuthInterceptor(
      tokenStorage: tokenStorage,
      dio: dio,
      onSessionExpired: () {
        // Handled through token clearing; auth state will sync
      },
    ),
  );

  return dio;
});
