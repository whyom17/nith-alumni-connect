import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:nith_alumni_connect/app/app.dart';
import 'package:nith_alumni_connect/core/storage/token_storage.dart';
import 'package:nith_alumni_connect/core/theme/theme_mode_provider.dart';
import 'package:nith_alumni_connect/features/auth/data/models/auth_dto.dart';
import 'package:nith_alumni_connect/features/auth/data/repositories/auth_repository.dart';
import 'package:nith_alumni_connect/features/auth/domain/models/auth_state.dart';
import 'package:nith_alumni_connect/features/auth/presentation/providers/auth_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class FakeAuthNotifier extends StateNotifier<AuthState>
    implements AuthNotifier {
  FakeAuthNotifier([AuthState? initial])
      : super(initial ?? const AuthState.unauthenticated());

  @override
  AuthRepository get authRepository => throw UnimplementedError();

  @override
  TokenStorage get tokenStorage => throw UnimplementedError();

  @override
  Future<void> initSession() async {}

  @override
  Future<AuthResponseDto> loginWithLinkedInCode({
    required String code,
    required String redirectUri,
    String? codeVerifier,
  }) async {
    const dummyUser = AuthUserDto(
      id: 'test-user',
      email: 'test@nith.ac.in',
      role: 'alumnus',
      status: 'pending',
      profileComplete: false,
    );
    state = const AuthState.authenticated(dummyUser);
    return const AuthResponseDto(
      accessToken: 'dummy-access',
      refreshToken: 'dummy-refresh',
      isNewUser: false,
      user: dummyUser,
    );
  }

  @override
  Future<void> logout() async {
    state = const AuthState.unauthenticated();
  }
}

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('Welcome screen renders and navigates to signup and login',
      (tester) async {
    final prefs = await SharedPreferences.getInstance();

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
          authNotifierProvider
              .overrideWith((ref) => FakeAuthNotifier(const AuthState.unauthenticated())),
        ],
        child: const AlumniConnectApp(),
      ),
    );

    await tester.pumpAndSettle();

    // Verify Welcome Screen elements
    expect(find.text('NITH ALUMNI CONNECT'), findsOneWidget);
    expect(find.text('Get Started'), findsOneWidget);
    expect(find.text('Log In'), findsOneWidget);

    // Tap Get Started -> should navigate to Create Account
    await tester.tap(find.text('Get Started'));
    await tester.pumpAndSettle();

    expect(find.text('Create Account'), findsOneWidget);
    expect(find.text('Sign up with LinkedIn'), findsOneWidget);

    // Tap "Log in" link -> navigates to Log In screen
    await tester.tap(find.text('Log in'));
    await tester.pumpAndSettle();

    expect(find.text('Log In'), findsOneWidget);
    expect(find.text('Log in with LinkedIn'), findsOneWidget);
  });
}
