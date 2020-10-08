import java.security.MessageDigest;

public final class I2b2PasswordHash {
    public static void main(String[] args) throws Exception {
        if (args.length != 1) {
            System.err.println("Usage: java I2b2PasswordHash <plaintext password>");
            System.exit(-1);
        }

        final MessageDigest md5 = MessageDigest.getInstance("MD5");
        md5.update(args[0].getBytes());

        byte[] digest = md5.digest();
        final StringBuffer buf = new StringBuffer();
        for (int i = 0; i < digest.length; i++) {
            buf.append(Integer.toHexString((int) digest[i] & 0x00FF));
        }

        System.out.println(buf.toString());
    }
}
