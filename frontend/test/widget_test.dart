import 'package:flutter_test/flutter_test.dart';
import 'package:show_love/app.dart';

void main() {
  testWidgets('ShowLoveApp renders feed page', (WidgetTester tester) async {
    await tester.pumpWidget(const ShowLoveApp());
    await tester.pumpAndSettle();
    // Initial route is /feed, should show the app title
    expect(find.text('显出爱心'), findsWidgets);
  });
}
