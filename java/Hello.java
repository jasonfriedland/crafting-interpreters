import java.lang.System.Logger.Level;

// Hello world class.
class Hello {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
        System.Logger log = System.getLogger("toot");
        log.log(Level.INFO, "done");
    }
}