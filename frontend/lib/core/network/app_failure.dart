/// Represents a typed domain error mapped from API or network failures.
sealed class AppFailure implements Exception {
  const AppFailure(this.message, {this.code});

  final String message;
  final String? code;

  @override
  String toString() => '$runtimeType(message: $message, code: $code)';
}

final class NetworkFailure extends AppFailure {
  const NetworkFailure([super.message = 'No internet connection.'])
      : super(code: 'network_error');
}

final class UnauthorizedFailure extends AppFailure {
  const UnauthorizedFailure([super.message = 'Please log in to continue.'])
      : super(code: 'unauthorized');
}

final class ForbiddenFailure extends AppFailure {
  const ForbiddenFailure([super.message = 'Access forbidden.'])
      : super(code: 'forbidden');
}

final class NotFoundFailure extends AppFailure {
  const NotFoundFailure([super.message = 'Resource not found.'])
      : super(code: 'not_found');
}

final class ValidationFailure extends AppFailure {
  const ValidationFailure(super.message, {this.fields})
      : super(code: 'validation_error');

  final Map<String, String>? fields;
}

final class ServerFailure extends AppFailure {
  const ServerFailure([super.message = 'Server encountered an error.'])
      : super(code: 'internal_error');
}

final class UnknownFailure extends AppFailure {
  const UnknownFailure([super.message = 'An unexpected error occurred.'])
      : super(code: 'unknown');
}
