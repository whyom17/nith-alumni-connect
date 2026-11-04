import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:nith_alumni_connect/app/app.dart';
import 'package:nith_alumni_connect/core/network/health_repository.dart';
import 'package:nith_alumni_connect/core/theme/theme_mode_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('App renders 4 navigation tabs and switches destinations',
      (tester) async {
    final prefs = await SharedPreferences.getInstance();

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
          healthCheckProvider.overrideWith((ref) => Future.value('ok')),
        ],
        child: const AlumniConnectApp(),
      ),
    );

    // Initial pump and settle
    await tester.pumpAndSettle();

    // Verify Home tab is selected by default
    expect(find.text('Alumni Connect'), findsOneWidget);
    expect(find.text('Backend Connection (M0)'), findsOneWidget);
    expect(find.text('Status: ok (Connected to /healthz)'), findsOneWidget);

    // Verify all 4 tabs exist in NavigationBar
    expect(find.text('Home'), findsOneWidget);
    expect(find.text('Groups'), findsOneWidget);
    expect(find.text('Search'), findsOneWidget);
    expect(find.text('Profile'), findsOneWidget);

    // Tap on Groups tab
    await tester.tap(find.text('Groups'));
    await tester.pumpAndSettle();

    expect(find.text('Groups & Communities'), findsOneWidget);
    expect(find.text('Communities Coming Soon'), findsOneWidget);

    // Tap on Search tab
    await tester.tap(find.text('Search'));
    await tester.pumpAndSettle();

    expect(find.text('Search Alumni'), findsOneWidget);
    expect(find.text('Alumni Directory'), findsOneWidget);

    // Tap on Profile tab
    await tester.tap(find.text('Profile'));
    await tester.pumpAndSettle();

    expect(find.text('My Profile'), findsOneWidget);
    expect(find.text('Alumni Profile'), findsOneWidget);

    // Tap back to Home
    await tester.tap(find.text('Home'));
    await tester.pumpAndSettle();
    expect(find.text('Alumni Connect'), findsOneWidget);
  });

  testWidgets('Drawer opens and contains navigation links and theme toggle',
      (tester) async {
    final prefs = await SharedPreferences.getInstance();

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          sharedPreferencesProvider.overrideWithValue(prefs),
          healthCheckProvider.overrideWith((ref) => Future.value('ok')),
        ],
        child: const AlumniConnectApp(),
      ),
    );

    await tester.pumpAndSettle();

    // Tap hamburger menu icon on Home screen
    await tester.tap(find.byTooltip('Open navigation menu'));
    await tester.pumpAndSettle();

    // Drawer should now be visible
    expect(find.text('NITH Alumnus'), findsOneWidget);
    expect(find.text('Class of 2024 | CSE'), findsOneWidget);
    expect(find.text('Saved Posts'), findsOneWidget);
    expect(find.text('Settings & Privacy'), findsOneWidget);
    expect(find.text('Log Out'), findsOneWidget);

    // Tap Saved Posts
    await tester.tap(find.text('Saved Posts'));
    await tester.pumpAndSettle();

    // Pushed route: Saved Posts screen
    expect(find.text('Saved Announcements'), findsOneWidget);
  });
}
