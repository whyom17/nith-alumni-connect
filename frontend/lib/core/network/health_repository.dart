import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'api_client.dart';

/// Repository responsible for querying backend liveness and readiness probes.
class HealthRepository {
  const HealthRepository(this._dio);

  final Dio _dio;

  /// Calls GET /healthz and returns the liveness status string.
  Future<String> checkHealth() async {
    try {
      final response = await _dio.get<Map<String, dynamic>>('/healthz');
      final data = response.data;
      return data?['status'] as String? ?? 'unknown';
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }

  /// Calls GET /readyz and returns the readiness status string.
  Future<String> checkReadiness() async {
    try {
      final response = await _dio.get<Map<String, dynamic>>('/readyz');
      final data = response.data;
      return data?['status'] as String? ?? 'unknown';
    } on DioException catch (e) {
      throw mapDioException(e);
    }
  }
}

/// Provides the [HealthRepository] singleton.
final healthRepositoryProvider = Provider<HealthRepository>((ref) {
  return HealthRepository(ref.watch(dioProvider));
});

/// Evaluates and provides the result of calling /healthz.
final healthCheckProvider = FutureProvider.autoDispose<String>((ref) {
  return ref.watch(healthRepositoryProvider).checkHealth();
});
