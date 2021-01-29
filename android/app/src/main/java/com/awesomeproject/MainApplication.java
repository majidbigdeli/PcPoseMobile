package com.awesomeproject;

import android.app.Application;
import android.content.Context;
import android.renderscript.ScriptIntrinsicYuvToRGB;
import android.util.Log;

import com.facebook.react.PackageList;
import com.facebook.react.ReactApplication;
import com.facebook.react.ReactInstanceManager;
import com.facebook.react.ReactNativeHost;
import com.facebook.react.ReactPackage;
import com.facebook.react.bridge.ReactContext;
import com.facebook.react.modules.core.DeviceEventManagerModule;
import com.facebook.react.shell.MainReactPackage;
import com.facebook.soloader.SoLoader;

import java.lang.reflect.InvocationTargetException;
import java.net.InetAddress;
import java.net.Socket;
import java.util.List;
import java.net.ServerSocket;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.ServerSocket;
import java.net.Socket;

import android.app.Activity;
import android.os.Bundle;
import android.os.Handler;
import android.widget.TextView;


import hello.Hello;
import hello.JavaCallback;
import hello.Person;


public class MainApplication extends Application implements ReactApplication {

    private final ReactNativeHost mReactNativeHost =
            new ReactNativeHost(this) {
                @Override
                public boolean getUseDeveloperSupport() {
                    return BuildConfig.DEBUG;
                }

                @Override
                protected List<ReactPackage> getPackages() {
                    @SuppressWarnings("UnnecessaryLocalVariable")
                    List<ReactPackage> packages = new PackageList(this).getPackages();
                    packages.add(new HelloPackage());
                    // Packages that cannot be autolinked yet can be added manually here, for example:
                    // packages.add(new MyReactNativePackage());
                    return packages;
                }

                @Override
                protected String getJSMainModuleName() {
                    return "index";
                }
            };

    private ServerSocket serverSocket;

    Handler updateConversationHandler;

    Thread serverThread = null;

    static JavaCallbackImp gocb;


    public static final int SERVERPORT = 6000;

    @Override
    public ReactNativeHost getReactNativeHost() {
        return mReactNativeHost;
    }

    @Override
    public void onCreate() {
        super.onCreate();


//    Thread t1 = new Thread(() ->{
//
//        Socket socket = null;
//        try {
//            InetAddress addr = InetAddress.getByName("10.0.2.15");
//
//            serverSocket = new ServerSocket(SERVERPORT, 50, addr);
//        } catch (IOException e) {
//            e.printStackTrace();
//        }
//        while (!Thread.currentThread().isInterrupted()) {
//
//            try {
//
//                socket = serverSocket.accept();
//
//                CommunicationThread commThread = new CommunicationThread(socket);
//                new Thread(commThread).start();
//
//            } catch (IOException e) {
//                e.printStackTrace();
//            }
//        }
//
//
//    });
//
//    t1.start();

//        new Handler().postDelayed(new Runnable() {
//            @Override
//            public void run() {
//                ReactContext reactContext = getReactNativeHost().getReactInstanceManager().getCurrentReactContext();
//            }
//        }, 1000);

        gocb = new JavaCallbackImp();
        Hello.registerJavaCallback(gocb);
        Hello.start();


        SoLoader.init(this, /* native exopackage */ false);
        initializeFlipper(this, getReactNativeHost().getReactInstanceManager());
    }


    /**
     * Loads Flipper in React Native templates. Call this in the onCreate method with something like
     * initializeFlipper(this, getReactNativeHost().getReactInstanceManager());
     *
     * @param context
     * @param reactInstanceManager
     */
    private static void initializeFlipper(
            Context context, ReactInstanceManager reactInstanceManager) {
        if (BuildConfig.DEBUG) {
            try {
        /*
         We use reflection here to pick up the class that initializes Flipper,
        since Flipper library is not available in release mode
        */
                Class<?> aClass = Class.forName("com.awesomeproject.ReactNativeFlipper");
                aClass
                        .getMethod("initializeFlipper", Context.class, ReactInstanceManager.class)
                        .invoke(null, context, reactInstanceManager);
            } catch (ClassNotFoundException e) {
                e.printStackTrace();
            } catch (NoSuchMethodException e) {
                e.printStackTrace();
            } catch (IllegalAccessException e) {
                e.printStackTrace();
            } catch (InvocationTargetException e) {
                e.printStackTrace();
            }
        }
    }


    class CommunicationThread implements Runnable {

        private Socket clientSocket;

        private BufferedReader input;

        public CommunicationThread(Socket clientSocket) {

            this.clientSocket = clientSocket;

            try {

                this.input = new BufferedReader(new InputStreamReader(this.clientSocket.getInputStream()));

            } catch (IOException e) {
                e.printStackTrace();
            }
        }

        public void run() {

            while (!Thread.currentThread().isInterrupted()) {

                try {

                    String read = input.readLine();

                    Log.d("FF", read);
//                    updateConversationHandler.post(new updateUIThread(read));

                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }

    }

    class JavaCallbackImp implements JavaCallback {

        public JavaCallbackImp() {

        }


        @Override
        public void callFromGo(String s) {
            Log.d("DD", s);
            try {
                getReactNativeHost().getReactInstanceManager().getCurrentReactContext()
                        .getJSModule(DeviceEventManagerModule.RCTDeviceEventEmitter.class)
                        .emit("customEventName", s);

            } catch (Exception e) {
                Log.e("ReactNative", "Caught Exception: " + e.getMessage());
            }
        }
    }

}


