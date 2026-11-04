import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:nith_alumni_connect/core/network/api_client.dart';
import 'package:nith_alumni_connect/core/network/app_failure.dart';

void main() {
  group('mapDioException', () {
    test('maps connection timeout to NetworkFailure', () {
      final dioError = DioException(
        requestOptions: RequestOptions(path: '/healthz'),
        type: DioExceptionType.connectionTimeout,
      );

      final failure = mapDioException(dioError);
      expect(failure, isA<NetworkFailure>());
    });

    test('maps 400 Bad Request to ValidationFailure', () {
      final dioError = DioException(
        requestOptions: RequestOptions(path: '/test'),
        response: Response(
          requestOptions: RequestOptions(path: '/test'),
          statusCode: 400,
          data: {
            'error': {
              'code': 'validation_error',
              'message': 'Invalid input',
              'fields': {'batch': 'must be a 4-digit year'},
            },
          },
        ),
      );

      final failure = mapDioException(dioError);
      expect(failure, isA<ValidationFailure>());
      final valFailure = failure as ValidationFailure;
      expect(valFailure.message, 'Invalid input');
      expect(valFailure.fields?['batch'], 'must be a 4-digit year');
    });

    test('maps 401 Unauthorized to UnauthorizedFailure', () {
      final dioError = DioException(
        requestOptions: RequestOptions(path: '/v1/me'),
        response: Response(
          requestOptions: RequestOptions(path: '/v1/me'),
          statusCode: 401,
          data: {
            'error': {
              'code': 'unauthorized',
              'message': 'Missing or expired token',
            },
          },
        ),
      );

      final failure = mapDioException(dioError);
      expect(failure, isA<UnauthorizedFailure>());
    });

    test('maps 503 to ServerFailure', () {
      final dioError = DioException(
        requestOptions: RequestOptions(path: '/readyz'),
        response: Response(
          requestOptions: RequestOptions(path: '/readyz'),
          statusCode: 503,
          data: {
            'error': {
              'code': 'service_unavailable',
              'message': 'database is not ready',
            },
          },
        ),
      );

      final failure = mapDioException(dioError);
      expect(failure, isA<ServerFailure>());
    });
  });
}
