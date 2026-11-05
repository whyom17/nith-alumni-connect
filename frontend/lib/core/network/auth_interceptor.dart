import 'package:dio/dio.dart';
import '../storage/token_storage.dart';

/// Interceptor that attaches JWT access tokens to outgoing requests and handles
/// automatic token rotation on 401 Unauthorized responses.
class AuthInterceptor extends QueuedInterceptor {
  AuthInterceptor({
    required this.tokenStorage,
    required this.dio,
    required this.onSessionExpired,
  });

  final TokenStorage tokenStorage;
  final Dio dio;
  final void Function() onSessionExpired;

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Skip attaching bearer tokens for unauthenticated auth endpoints
    if (options.path.startsWith('/v1/auth/linkedin') ||
        options.path.startsWith('/v1/auth/refresh')) {
      return handler.next(options);
    }

    final token = await tokenStorage.getAccessToken();
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    return handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    // Only attempt refresh on 401 errors for non-auth paths
    if (err.response?.statusCode == 401 &&
        !err.requestOptions.path.startsWith('/v1/auth/')) {
      final refreshToken = await tokenStorage.getRefreshToken();
      if (refreshToken == null || refreshToken.isEmpty) {
        onSessionExpired();
        return handler.next(err);
      }

      try {
        // Attempt token rotation directly using a clean Dio instance to avoid re-triggering interceptor
        final refreshDio = Dio(
          BaseOptions(
            baseUrl: dio.options.baseUrl,
            connectTimeout: const Duration(seconds: 10),
            receiveTimeout: const Duration(seconds: 10),
          ),
        );

        final response = await refreshDio.post<Map<String, dynamic>>(
          '/v1/auth/refresh',
          data: {'refresh_token': refreshToken},
        );

        final data = response.data;
        if (data != null && data['access_token'] is String) {
          final newAccess = data['access_token'] as String;
          final newRefresh = data['refresh_token'] as String;

          await tokenStorage.saveTokens(
            accessToken: newAccess,
            refreshToken: newRefresh,
          );

          // Retry the original request with new access token
          final requestOptions = err.requestOptions;
          requestOptions.headers['Authorization'] = 'Bearer $newAccess';

          final clonedResponse = await dio.fetch(requestOptions);
          return handler.resolve(clonedResponse);
        }
      } catch (_) {
        // Refresh failed (token revoked/expired) - clear tokens and trigger sign out
        await tokenStorage.clearTokens();
        onSessionExpired();
      }
    }

    return handler.next(err);
  }
}
