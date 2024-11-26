import java.util.Scanner;

public class MathUtilities {
    
    // Function to compute the Fibonacci sequence
    public static void fibonacci(int terms) {
        System.out.println("Fibonacci sequence with " + terms + " terms:");
        int a = 0, b = 1;
        for (int i = 0; i < terms; i++) {
            System.out.print(a + " ");
            int next = a + b;
            a = b;
            b = next;
        }
        System.out.println();
    }

    // Function to compute the GCD of two numbers
    public static int gcd(int m, int n) {
        while (n != 0) {
            int temp = n;
            n = m % n;
            m = temp;
        }
        return m;
    }

    // Function to compute the roots of a quadratic equation
    public static void quadraticRoots(double a, double b, double c) {
        if (a == 0) {
            System.out.println("This is not a quadratic equation (a cannot be 0).");
            return;
        }

        double discriminant = b * b - 4 * a * c;

        if (discriminant > 0) {
            double root1 = (-b + Math.sqrt(discriminant)) / (2 * a);
            double root2 = (-b - Math.sqrt(discriminant)) / (2 * a);
            System.out.println("The roots are real and distinct:");
            System.out.println("Root 1: " + root1);
            System.out.println("Root 2: " + root2);
        } else if (discriminant == 0) {
            double root = -b / (2 * a);
            System.out.println("The roots are real and equal:");
            System.out.println("Root: " + root);
        } else {
            double realPart = -b / (2 * a);
            double imaginaryPart = Math.sqrt(-discriminant) / (2 * a);
            System.out.println("The roots are complex:");
            System.out.println("Root 1: " + realPart + " + " + imaginaryPart + "i");
            System.out.println("Root 2: " + realPart + " - " + imaginaryPart + "i");
        }
    }

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // Fibonacci sequence
        System.out.print("Enter the number of terms for the Fibonacci sequence: ");
        int terms = scanner.nextInt();
        fibonacci(terms);

        // GCD
        System.out.print("Enter the first number (m): ");
        int m = scanner.nextInt();
        System.out.print("Enter the second number (n): ");
        int n = scanner.nextInt();
        System.out.println("The GCD of " + m + " and " + n + " is: " + gcd(m, n));

        // Quadratic equation
        System.out.print("Enter coefficient a: ");
        double a = scanner.nextDouble();
        System.out.print("Enter coefficient b: ");
        double b = scanner.nextDouble();
        System.out.print("Enter coefficient c: ");
        double c = scanner.nextDouble();
        quadraticRoots(a, b, c);

        scanner.close();
    }
}
